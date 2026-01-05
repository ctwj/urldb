package jsvm

import (
	"sync"

	"github.com/dop251/goja"
)

type poolItem struct {
	mux  sync.Mutex
	busy bool
	vm   *goja.Runtime
}

type vmsPool struct {
	mux     sync.RWMutex
	factory func() *goja.Runtime
	items   []*poolItem
}

// newPool creates a new pool with pre-warmed vms generated from the specified factory.
func newPool(size int, factory func() *goja.Runtime) *vmsPool {
	pool := &vmsPool{
		factory: factory,
		items:   make([]*poolItem, size),
	}

	for i := 0; i < size; i++ {
		vm := pool.factory()
		pool.items[i] = &poolItem{vm: vm}
	}

	return pool
}

// run executes "call" with a vm created from the pool
// (either from the buffer or a new one if all buffered vms are busy)
func (p *vmsPool) run(call func(vm *goja.Runtime) error) error {
	p.mux.RLock()

	// try to find a free item
	var freeItem *poolItem
	for _, item := range p.items {
		item.mux.Lock()
		if item.busy {
			item.mux.Unlock()
			continue
		}
		item.busy = true
		item.mux.Unlock()
		freeItem = item
		break
	}

	p.mux.RUnlock()

	// create a new one-off item if of all of the pool items are currently busy
	//
	// note: if turned out not efficient we may change this in the future
	// by adding the created item in the pool with some timer for removal
	if freeItem == nil {
		return call(p.factory())
	}

	execErr := call(freeItem.vm)

	// "free" the vm
	freeItem.mux.Lock()
	freeItem.busy = false
	freeItem.mux.Unlock()

	return execErr
}

// Get gets a vm from the pool (similar to run but returns the vm directly)
func (p *vmsPool) Get() *goja.Runtime {
	p.mux.RLock()

	// try to find a free item
	var freeItem *poolItem
	for _, item := range p.items {
		item.mux.Lock()
		if item.busy {
			item.mux.Unlock()
			continue
		}
		item.busy = true
		item.mux.Unlock()
		freeItem = item
		break
	}

	p.mux.RUnlock()

	// create a new one-off item if all pool items are currently busy
	if freeItem == nil {
		return p.factory()
	}

	return freeItem.vm
}

// Put returns a vm to the pool (note: this is a simplified implementation)
func (p *vmsPool) Put(vm *goja.Runtime) {
	// Find the item in the pool and mark it as not busy
	for _, item := range p.items {
		if item.vm == vm {
			item.mux.Lock()
			item.busy = false
			item.mux.Unlock()
			return
		}
	}
	// If not found in pool, it was a temporary vm, nothing to do
}
