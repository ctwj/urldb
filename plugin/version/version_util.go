package version

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Semantic versioning helper functions

// IsValidSemanticVersion 检查是否为有效的语义化版本
func IsValidSemanticVersion(version string) bool {
	// 语义化版本格式: MAJOR.MINOR.PATCH[-PRERELEASE][+BUILDINFO]
	// 例如: 1.0.0, 1.0.0-alpha, 1.0.0+build.1
	pattern := `^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`

	re := regexp.MustCompile(pattern)
	return re.MatchString(version)
}

// ParseSemanticVersion 解析语义化版本
func ParseSemanticVersion(version string) (major, minor, patch int, prerelease, build string, err error) {
	parts := strings.Split(version, "+")
	if len(parts) > 1 {
		build = parts[1]
		version = parts[0]
	}

	parts = strings.Split(version, "-")
	if len(parts) > 1 {
		prerelease = parts[1]
		version = parts[0]
	}

	numbers := strings.Split(version, ".")
	if len(numbers) != 3 {
		err = fmt.Errorf("invalid version format: %s", version)
		return
	}

	major, err = strconv.Atoi(numbers[0])
	if err != nil {
		return
	}

	minor, err = strconv.Atoi(numbers[1])
	if err != nil {
		return
	}

	patch, err = strconv.Atoi(numbers[2])
	if err != nil {
		return
	}

	return
}

// CompareSemanticVersions 比较两个语义化版本
func CompareSemanticVersions(v1, v2 string) int {
	// 如果版本相等，直接返回0
	if v1 == v2 {
		return 0
	}

	major1, minor1, patch1, prerelease1, _, err1 := ParseSemanticVersion(v1)
	major2, minor2, patch2, prerelease2, _, err2 := ParseSemanticVersion(v2)

	// 如果解析失败，使用字符串比较作为备选方案
	if err1 != nil || err2 != nil {
		return compareStringVersions(v1, v2)
	}

	// 比较主版本号
	if major1 > major2 {
		return 1
	} else if major1 < major2 {
		return -1
	}

	// 比较次版本号
	if minor1 > minor2 {
		return 1
	} else if minor1 < minor2 {
		return -1
	}

	// 比较修订号
	if patch1 > patch2 {
		return 1
	} else if patch1 < patch2 {
		return -1
	}

	// 如果都是稳定版本，版本相同
	if prerelease1 == "" && prerelease2 == "" {
		return 0
	}

	// 预发布版本总是比稳定版本低
	if prerelease1 == "" && prerelease2 != "" {
		return 1
	}
	if prerelease1 != "" && prerelease2 == "" {
		return -1
	}

	// 比较预发布版本
	return comparePreReleaseVersions(prerelease1, prerelease2)
}

// compareStringVersions 字符串版本比较（备选方案）
func compareStringVersions(v1, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var num1, num2 int

		if i < len(parts1) {
			// 尝试转换为数字，如果失败则按字符串处理
			if n, err := strconv.Atoi(parts1[i]); err == nil {
				num1 = n
			} else {
				// 如果包含字母，按字符串比较
				if i < len(parts2) {
					if n2, err2 := strconv.Atoi(parts2[i]); err2 == nil {
						// 数字版本大于字符串版本
						return 1
					} else {
						// 字符串比较
						if parts1[i] > parts2[i] {
							return 1
						} else if parts1[i] < parts2[i] {
							return -1
						}
					}
				} else {
					return 1
				}
			}
		}

		if i < len(parts2) {
			if n, err := strconv.Atoi(parts2[i]); err == nil {
				num2 = n
			} else {
				// 如果v1已经处理了数字而v2是字符串
				if i < len(parts1) {
					if _, err2 := strconv.Atoi(parts1[i]); err2 == nil {
						return -1
					} else {
						if parts1[i] > parts2[i] {
							return 1
						} else if parts1[i] < parts2[i] {
							return -1
						}
					}
				} else {
					return -1
				}
			}
		}

		if num1 > num2 {
			return 1
		} else if num1 < num2 {
			return -1
		}
	}

	return 0
}

// comparePreReleaseVersions 比较预发布版本
func comparePreReleaseVersions(prerelease1, prerelease2 string) int {
	if prerelease1 == prerelease2 {
		return 0
	}

	// 如果其中一个预发布版本为空，非空的版本更小
	if prerelease1 == "" {
		return 1
	}
	if prerelease2 == "" {
		return -1
	}

	// 按标识符分割并比较
	parts1 := strings.Split(prerelease1, ".")
	parts2 := strings.Split(prerelease2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var part1, part2 string

		if i < len(parts1) {
			part1 = parts1[i]
		}
		if i < len(parts2) {
			part2 = parts2[i]
		}

		// 如果其中一个为空，非空的更大
		if part1 == "" {
			return -1
		}
		if part2 == "" {
			return 1
		}

		// 检查是否为数字
		num1, err1 := strconv.Atoi(part1)
		num2, err2 := strconv.Atoi(part2)

		if err1 == nil && err2 == nil {
			// 都是数字，数值比较
			if num1 > num2 {
				return 1
			} else if num1 < num2 {
				return -1
			}
		} else if err1 == nil {
			// 只有part1是数字，数字小于字符串
			return -1
		} else if err2 == nil {
			// 只有part2是数字，字符串大于数字
			return 1
		} else {
			// 都是字符串，字典序比较
			if part1 > part2 {
				return 1
			} else if part1 < part2 {
				return -1
			}
		}
	}

	return 0
}

// GetNextVersion 获取下一个版本号
func GetNextVersion(currentVersion string, bumpType string) (string, error) {
	major, minor, patch, prerelease, build, err := ParseSemanticVersion(currentVersion)
	if err != nil {
		return "", fmt.Errorf("invalid version format: %v", err)
	}

	switch bumpType {
	case "major":
		major++
		minor = 0
		patch = 0
		prerelease = ""
	case "minor":
		minor++
		patch = 0
		prerelease = ""
	case "patch":
		patch++
		prerelease = ""
	case "prerelease":
		// 如果当前是稳定版本，添加预发布标识符
		if prerelease == "" {
			prerelease = "alpha"
		} else {
			// 如果已经是预发布版本，递增数字部分
			if num, err := strconv.Atoi(prerelease); err == nil {
				prerelease = fmt.Sprintf("%d", num+1)
			} else {
				prerelease = fmt.Sprintf("%s.1", prerelease)
			}
		}
	default:
		return "", fmt.Errorf("invalid bump type: %s", bumpType)
	}

	nextVersion := fmt.Sprintf("%d.%d.%d", major, minor, patch)
	if prerelease != "" {
		nextVersion = fmt.Sprintf("%s-%s", nextVersion, prerelease)
	}
	if build != "" {
		nextVersion = fmt.Sprintf("%s+%s", nextVersion, build)
	}

	return nextVersion, nil
}

// IsStableVersion 检查是否为稳定版本
func IsStableVersion(version string) bool {
	// 稳定版本不包含预发布标识符
	parts := strings.Split(version, "-")
	return len(parts) == 1 || parts[1] == ""
}

// IsPrereleaseVersion 检查是否为预发布版本
func IsPrereleaseVersion(version string) bool {
	// 预发布版本包含预发布标识符
	parts := strings.Split(version, "-")
	return len(parts) > 1 && parts[1] != ""
}

// GetMajorVersion 获取主版本号
func GetMajorVersion(version string) (int, error) {
	major, _, _, _, _, err := ParseSemanticVersion(version)
	return major, err
}

// GetMinorVersion 获取次版本号
func GetMinorVersion(version string) (int, error) {
	_, minor, _, _, _, err := ParseSemanticVersion(version)
	return minor, err
}

// GetPatchVersion 获取修订版本号
func GetPatchVersion(version string) (int, error) {
	_, _, patch, _, _, err := ParseSemanticVersion(version)
	return patch, err
}