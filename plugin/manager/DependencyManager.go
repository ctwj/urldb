package manager

import (
	"fmt"
	"sort"
	"sync"

	"github.com/ctwj/urldb/plugin/types"
	"github.com/ctwj/urldb/utils"
)

// DependencyManager handles plugin dependency resolution and loading order
type DependencyManager struct {
	manager *Manager
	mutex   sync.RWMutex
}

// NewDependencyManager creates a new dependency manager
func NewDependencyManager(manager *Manager) *DependencyManager {
	return &DependencyManager{
		manager: manager,
	}
}

// DependencyGraph represents the dependency relationships between plugins
type DependencyGraph struct {
	adjacencyList map[string][]string // plugin -> dependencies
	reverseList   map[string][]string // plugin -> dependents of this plugin
}

// NewDependencyGraph creates a new dependency graph
func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		adjacencyList: make(map[string][]string),
		reverseList:   make(map[string][]string),
	}
}

// AddDependency adds a dependency relationship
func (dg *DependencyGraph) AddDependency(plugin, dependency string) {
	if _, exists := dg.adjacencyList[plugin]; !exists {
		dg.adjacencyList[plugin] = []string{}
	}
	if _, exists := dg.reverseList[dependency]; !exists {
		dg.reverseList[dependency] = []string{}
	}

	// Add dependency
	dg.adjacencyList[plugin] = append(dg.adjacencyList[plugin], dependency)
	// Add reverse dependency (plugin depends on dependency, so dependency is needed by plugin)
	dg.reverseList[dependency] = append(dg.reverseList[dependency], plugin)
}

// GetDependencies returns all direct dependencies of a plugin
func (dg *DependencyGraph) GetDependencies(plugin string) []string {
	if deps, exists := dg.adjacencyList[plugin]; exists {
		return deps
	}
	return []string{}
}

// GetDependents returns all plugins that depend on the given plugin
func (dg *DependencyGraph) GetDependents(plugin string) []string {
	if deps, exists := dg.reverseList[plugin]; exists {
		return deps
	}
	return []string{}
}

// GetAllPlugins returns all plugins in the graph
func (dg *DependencyGraph) GetAllPlugins() []string {
	plugins := make([]string, 0)
	for plugin := range dg.adjacencyList {
		plugins = append(plugins, plugin)
	}
	// Also add plugins that are only depended on but don't have dependencies themselves
	for plugin := range dg.reverseList {
		found := false
		for _, p := range plugins {
			if p == plugin {
				found = true
				break
			}
		}
		if !found {
			plugins = append(plugins, plugin)
		}
	}
	return plugins
}

// ValidateDependencies validates that all dependencies exist
func (dm *DependencyManager) ValidateDependencies() error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	// Build dependency graph
	graph := dm.buildDependencyGraph()

	// Check that all dependencies exist
	for pluginName, plugin := range dm.manager.plugins {
		dependencies := plugin.Dependencies()
		for _, dep := range dependencies {
			if _, exists := dm.manager.plugins[dep]; !exists {
				return fmt.Errorf("plugin %s depends on %s, but %s is not registered", pluginName, dep, dep)
			}
		}
	}

	// Check for circular dependencies
	cycles := dm.findCircularDependencies(graph)
	if len(cycles) > 0 {
		return fmt.Errorf("circular dependencies detected: %v", cycles)
	}

	return nil
}

// CheckPluginDependencies checks if all dependencies for a specific plugin are satisfied
func (dm *DependencyManager) CheckPluginDependencies(pluginName string) (bool, []string, error) {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	plugin, exists := dm.manager.plugins[pluginName]
	if !exists {
		return false, nil, fmt.Errorf("plugin %s not found", pluginName)
	}

	dependencies := plugin.Dependencies()
	unresolved := []string{}

	for _, dep := range dependencies {
		depInstance, exists := dm.manager.instances[dep]
		if !exists {
			unresolved = append(unresolved, dep)
			continue
		}
		if depInstance.Status != types.StatusRunning {
			unresolved = append(unresolved, dep)
		}
	}

	return len(unresolved) == 0, unresolved, nil
}

// GetLoadOrder returns the correct order to load plugins based on dependencies
func (dm *DependencyManager) GetLoadOrder() ([]string, error) {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	graph := dm.buildDependencyGraph()

	// Check for circular dependencies
	cycles := dm.findCircularDependencies(graph)
	if len(cycles) > 0 {
		return nil, fmt.Errorf("circular dependencies detected: %v", cycles)
	}

	// Topological sort to get load order
	return dm.topologicalSort(graph), nil
}

// buildDependencyGraph builds a dependency graph from registered plugins
func (dm *DependencyManager) buildDependencyGraph() *DependencyGraph {
	graph := NewDependencyGraph()

	for pluginName, plugin := range dm.manager.plugins {
		dependencies := plugin.Dependencies()
		for _, dep := range dependencies {
			graph.AddDependency(pluginName, dep)
		}
	}

	return graph
}

// findCircularDependencies finds circular dependencies in the dependency graph
func (dm *DependencyManager) findCircularDependencies(graph *DependencyGraph) [][]string {
	var cycles [][]string

	// Use DFS to detect cycles
	visited := make(map[string]bool)
	recStack := make(map[string]bool)
	path := []string{}

	for _, plugin := range graph.GetAllPlugins() {
		if !visited[plugin] {
			cycleFound := dm.dfsForCycles(plugin, graph, visited, recStack, path, &cycles)
			if cycleFound {
				// If we found cycles, we can stop (or continue to find more)
				// For now, just continue to find all possible cycles
			}
		}
	}

	return cycles
}

// dfsForCycles performs DFS to detect cycles in the dependency graph
func (dm *DependencyManager) dfsForCycles(
	plugin string,
	graph *DependencyGraph,
	visited map[string]bool,
	recStack map[string]bool,
	path []string,
	cycles *[][]string,
) bool {
	visited[plugin] = true
	recStack[plugin] = true
	path = append(path, plugin)

	dependencies := graph.GetDependencies(plugin)
	for _, dep := range dependencies {
		if !visited[dep] {
			if dm.dfsForCycles(dep, graph, visited, recStack, path, cycles) {
				return true
			}
		} else if recStack[dep] {
			// Found a cycle
			// Find where this dependency first appears in the current path
			startIdx := -1
			for i, p := range path {
				if p == dep {
					startIdx = i
					break
				}
			}
			if startIdx != -1 {
				cycle := path[startIdx:]
				cycle = append(cycle, dep) // Add the dependency again to close the cycle
				*cycles = append(*cycles, cycle)
			}
			return true
		}
	}

	recStack[plugin] = false
	// Remove the last element from path
	if len(path) > 0 {
		path = path[:len(path)-1]
	}

	return false
}

// topologicalSort performs topological sort to get the correct load order
func (dm *DependencyManager) topologicalSort(graph *DependencyGraph) []string {
	var result []string
	visited := make(map[string]bool)
	temporary := make(map[string]bool)

	// Get all plugins
	allPlugins := graph.GetAllPlugins()

	// Sort plugins to ensure consistent order
	sort.Strings(allPlugins)

	var visit func(string) error
	visit = func(plugin string) error {
		if temporary[plugin] {
			return fmt.Errorf("circular dependency detected")
		}
		if !visited[plugin] {
			temporary[plugin] = true

			dependencies := graph.GetDependencies(plugin)
			for _, dep := range dependencies {
				if err := visit(dep); err != nil {
					return err
				}
			}

			temporary[plugin] = false
			visited[plugin] = true
			result = append(result, plugin)
		}
		return nil
	}

	for _, plugin := range allPlugins {
		if !visited[plugin] {
			if err := visit(plugin); err != nil {
				utils.Error("Error during topological sort: %v", err)
				return []string{} // Return empty slice if there's an error
			}
		}
	}

	return result
}

// GetDependencyInfo returns dependency information for a plugin
func (dm *DependencyManager) GetDependencyInfo(pluginName string) (*types.PluginInfo, error) {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	plugin, exists := dm.manager.plugins[pluginName]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", pluginName)
	}

	// Build dependency graph
	graph := dm.buildDependencyGraph()

	info := &types.PluginInfo{
		Name:         pluginName,
		Version:      plugin.Version(),
		Description:  plugin.Description(),
		Author:       plugin.Author(),
		Dependencies: graph.GetDependencies(pluginName),
	}

	// Get instance status if available
	if instance, exists := dm.manager.instances[pluginName]; exists {
		info.Status = instance.Status
		info.LastError = instance.LastError
		info.StartTime = instance.StartTime
		info.StopTime = instance.StopTime
		info.RestartCount = instance.RestartCount
		info.HealthScore = instance.HealthScore
	} else {
		info.Status = types.StatusRegistered
	}

	return info, nil
}

// CheckAllDependencies checks all plugin dependencies
func (dm *DependencyManager) CheckAllDependencies() map[string]map[string]bool {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	result := make(map[string]map[string]bool)

	for pluginName, plugin := range dm.manager.plugins {
		dependencies := plugin.Dependencies()
		depStatus := make(map[string]bool)

		for _, dep := range dependencies {
			depInstance, exists := dm.manager.instances[dep]
			depStatus[dep] = exists && (depInstance.Status == types.StatusRunning)
		}

		result[pluginName] = depStatus
	}

	return result
}

// GetDependents returns all plugins that depend on the given plugin
func (dm *DependencyManager) GetDependents(pluginName string) []string {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	graph := dm.buildDependencyGraph()
	return graph.GetDependents(pluginName)
}

// RemovePlugin removes a plugin from the dependency tracking
func (dm *DependencyManager) RemovePlugin(pluginName string) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	// The dependency graph will be rebuilt from the manager's plugin list,
	// so removing the plugin from the manager's list will effectively remove
	// it from the dependency tracking
	// No additional action needed here as the graph is built dynamically
}

// GetDependencies returns all direct dependencies of a plugin
func (dm *DependencyManager) GetDependencies(pluginName string) []string {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	graph := dm.buildDependencyGraph()
	return graph.GetDependencies(pluginName)
}