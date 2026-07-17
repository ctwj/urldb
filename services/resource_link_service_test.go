package services

import (
	"context"
	"testing"

	"github.com/ctwj/urldb/db/entity"
	"github.com/ctwj/urldb/db/repo"
)

// 本测试为纯单元测试：用 interface 嵌入（nil 默认）构造 fake 仓库与服务，
// 覆盖 ResolveWithCheck 决策树（FR-013~FR-017）的全部路由分支，不依赖 DB/外部网盘/PanCheck。
// 凡触达真实网盘 API 的「转存/分享成功」分支（需 live 账号）不在单测范围，
// 此处通过让 fake cksRepo 取账号失败，验证其「失败安全回退到原链」的降级路径。

// --- fakes ---

// fakeLinkCheck 按 URL 键返回受控结论（绕过 NormalizeURL，键须与 resource.URL/SaveURL 一致）
type fakeLinkCheck struct {
	LinkCheckService
	urls map[string]ResourceCheckResult
}

func (f *fakeLinkCheck) CheckURLs(_ context.Context, _ []string, _ bool) map[string]ResourceCheckResult {
	return f.urls
}

type fakePanRepo struct {
	repo.PanRepository
	pan *entity.Pan
}

func (f *fakePanRepo) FindByID(_ uint) (*entity.Pan, error) { return f.pan, nil }

type fakeConfigRepo struct {
	repo.SystemConfigRepository
	autoTransfer bool
}

func (f *fakeConfigRepo) GetConfigBool(key string) (bool, error) {
	if key == entity.ConfigKeyAutoTransferEnabled {
		return f.autoTransfer, nil
	}
	return false, nil
}

type fakeResourceRepo struct {
	repo.ResourceRepository
	updatedValid map[uint]bool
}

func (f *fakeResourceRepo) UpdateIsValid(id uint, valid bool) error {
	if f.updatedValid == nil {
		f.updatedValid = map[uint]bool{}
	}
	f.updatedValid[id] = valid
	return nil
}

func (f *fakeResourceRepo) UpdateFields(_ uint, _ map[string]interface{}) error { return nil }

// fakeCksRepo 让 PerformShare/PerformAutoTransfer 在「取账号」阶段失败，避免触达真实网盘 API
type fakeCksRepo struct {
	repo.CksRepository
}

func (f *fakeCksRepo) FindByID(_ uint) (*entity.Cks, error)        { return nil, errNoAccount }
func (f *fakeCksRepo) FindByPanID(_ uint) ([]entity.Cks, error)    { return nil, errNoAccount }

var errNoAccount = &fakeErr{"no account"}

type fakeErr struct{ s string }

func (e *fakeErr) Error() string { return e.s }

// --- helpers ---

func uptr(i uint) *uint { return &i }

func mkRes(id uint, url, saveURL string, isValid bool) *entity.Resource {
	r := &entity.Resource{URL: url, SaveURL: saveURL, IsValid: isValid}
	r.ID = id
	r.PanID = uptr(1)
	return r
}

var (
	panQuark   = &entity.Pan{ID: 1, Name: "quark", Remark: "夸克网盘"}
	panAlipan  = &entity.Pan{ID: 2, Name: "alipan", Remark: "阿里云盘"}
)

func newSvc(pan *entity.Pan, checkURLs map[string]ResourceCheckResult, autoTransfer bool, resRepo *fakeResourceRepo) ResourceLinkService {
	return NewResourceLinkService(
		&fakeCksRepo{},
		&fakePanRepo{pan: pan},
		&fakeConfigRepo{autoTransfer: autoTransfer},
		resRepo,
		&fakeLinkCheck{urls: checkURLs},
		nil, // meiliManager nil → ApplyValidityWriteback 跳过 Meili 同步
	)
}

func TestResolveWithCheck_DecisionTree(t *testing.T) {
	const (
		origURL = "https://pan.quark.cn/s/abc"
		saveURL = "https://pan.quark.cn/s/xyz"
	)
	validRes := ResourceCheckResult{Status: "valid", DetectionMethod: "pancheck"}
	invalidRes := ResourceCheckResult{Status: "invalid", DetectionMethod: "pancheck"}

	tests := []struct {
		name            string
		pan             *entity.Pan
		resource        *entity.Resource
		checkURLs       map[string]ResourceCheckResult
		autoTransfer    bool
		wantType        string
		wantInvalid     bool
		wantURL         string
		wantNoteEmpty   bool
		wantValidWrite  *bool // nil=期望无 is_valid 回写
	}{
		{
			name: "非转存平台·原始有效→原链",
			pan: panAlipan, resource: mkRes(1, origURL, "", true),
			checkURLs: map[string]ResourceCheckResult{origURL: validRes},
			wantType: "original", wantURL: origURL, wantNoteEmpty: true,
		},
		{
			name: "非转存平台·原始失效→Invalid+回写false",
			pan: panAlipan, resource: mkRes(2, origURL, "", true),
			checkURLs: map[string]ResourceCheckResult{origURL: invalidRes},
			wantType: "invalid", wantInvalid: true, wantNoteEmpty: true, wantValidWrite: bptr(false),
		},
		{
			name: "quark·saveUrl有效→transferred(直返saveUrl)",
			pan: panQuark, resource: mkRes(3, origURL, saveURL, true),
			checkURLs: map[string]ResourceCheckResult{origURL: validRes, saveURL: validRes},
			autoTransfer: true,
			wantType: "transferred", wantURL: saveURL, wantNoteEmpty: true,
		},
		{
			name: "quark·saveUrl未确定(PanCheck降级)→transferred(复用saveUrl)",
			pan: panQuark, resource: mkRes(4, origURL, saveURL, true),
			checkURLs: map[string]ResourceCheckResult{}, // 空=disabled，全 undetermined
			wantType: "transferred", wantURL: saveURL, wantNoteEmpty: true,
		},
		{
			name: "quark·saveUrl失效→分享失败+转存失败→original(静默回退,无错误Note)",
			pan: panQuark, resource: func() *entity.Resource {
				r := mkRes(5, origURL, saveURL, true)
				r.Fid = "fid-5"
				r.CkID = uptr(9)
				return r
			}(),
			checkURLs: map[string]ResourceCheckResult{origURL: validRes, saveURL: invalidRes},
			autoTransfer: true,
			wantType: "original", wantURL: origURL, wantNoteEmpty: true,
		},
		{
			name: "quark·无saveUrl·原始有效·自动转存关→original",
			pan: panQuark, resource: mkRes(6, origURL, "", true),
			checkURLs: map[string]ResourceCheckResult{origURL: validRes},
			autoTransfer: false,
			wantType: "original", wantURL: origURL, wantNoteEmpty: true,
		},
		{
			name: "quark·无saveUrl·原始失效→Invalid+回写false",
			pan: panQuark, resource: mkRes(7, origURL, "", true),
			checkURLs: map[string]ResourceCheckResult{origURL: invalidRes},
			wantType: "invalid", wantInvalid: true, wantNoteEmpty: true, wantValidWrite: bptr(false),
		},
		{
			name: "alipan·原始有效·IsValid=false→回写true+原链",
			pan: panAlipan, resource: mkRes(8, origURL, "", false),
			checkURLs: map[string]ResourceCheckResult{origURL: validRes},
			wantType: "original", wantURL: origURL, wantNoteEmpty: true, wantValidWrite: bptr(true),
		},
		{
			name: "quark·无saveUrl·原始未确定(降级)·转存关→original",
			pan: panQuark, resource: mkRes(9, origURL, "", true),
			checkURLs: map[string]ResourceCheckResult{},
			autoTransfer: false,
			wantType: "original", wantURL: origURL, wantNoteEmpty: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resRepo := &fakeResourceRepo{}
			svc := newSvc(tc.pan, tc.checkURLs, tc.autoTransfer, resRepo)

			got, err := svc.ResolveWithCheck(context.Background(), tc.resource)
			if err != nil {
				t.Fatalf("ResolveWithCheck 返回错误: %v", err)
			}
			if got.Type != tc.wantType {
				t.Errorf("Type = %q, want %q", got.Type, tc.wantType)
			}
			if got.Invalid != tc.wantInvalid {
				t.Errorf("Invalid = %v, want %v", got.Invalid, tc.wantInvalid)
			}
			if got.URL != tc.wantURL {
				t.Errorf("URL = %q, want %q", got.URL, tc.wantURL)
			}
			if (got.Note == "") != tc.wantNoteEmpty {
				t.Errorf("Note 空=%v, want 空=%v (Note=%q)", got.Note == "", tc.wantNoteEmpty, got.Note)
			}

			// is_valid 回写校验
			if tc.wantValidWrite != nil {
				v, ok := resRepo.updatedValid[tc.resource.ID]
				if !ok {
					t.Fatalf("期望回写 is_valid=%v，但未发生回写", *tc.wantValidWrite)
				}
				if v != *tc.wantValidWrite {
					t.Errorf("回写 is_valid=%v, want %v", v, *tc.wantValidWrite)
				}
			}
		})
	}
}

func bptr(b bool) *bool { return &b }

// TestPerformShare_Guards 验证 PerformShare 的快速失败守卫：
// 缺少 ck_id/fid 或取账号失败时立即返回失败，绝不触达真实网盘 API。
func TestPerformShare_Guards(t *testing.T) {
	resRepo := &fakeResourceRepo{}
	cksRepo := &fakeCksRepo{}

	tests := []struct {
		name       string
		resource   *entity.Resource
		wantSuccess bool
	}{
		{"resource 为空", nil, false},
		{"缺 ck_id", &entity.Resource{ID: 1, Fid: "fid"}, false},
		{"缺 fid", &entity.Resource{ID: 1, CkID: uptr(2)}, false},
		{"ck_id+fid 齐全但取账号失败(不触达网盘)", func() *entity.Resource {
			r := &entity.Resource{ID: 1, URL: "https://pan.quark.cn/s/x", Fid: "fid", SaveURL: "https://old"}
			r.CkID = uptr(2)
			r.PanID = uptr(1)
			return r
		}(), false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := PerformShare(cksRepo, resRepo, tc.resource)
			if got.Success != tc.wantSuccess {
				t.Errorf("Success = %v, want %v (ErrorMsg=%q)", got.Success, tc.wantSuccess, got.ErrorMsg)
			}
			if !tc.wantSuccess && got.ErrorMsg == "" {
				t.Errorf("失败时应有 ErrorMsg，实际为空")
			}
		})
	}
}

