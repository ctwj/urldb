package cmdplugin

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ctwj/urldb/utils"
	"github.com/spf13/cobra"
)

// pluginCmd 插件管理命令
var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "插件系统管理命令",
	Long:  `管理 URLDB 插件系统，包括创建模板、列出插件、验证文件等功能`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// GetPluginCmd 获取插件命令
func GetPluginCmd() *cobra.Command {
	return pluginCmd
}

// createCmd 创建插件命令
var createCmd = &cobra.Command{
	Use:   "create [name] [type]",
	Short: "创建新的插件模板",
	Long: `创建一个新的插件模板文件

支持的类型:
  hook     - 创建钩子插件模板
  migration - 创建数据库迁移模板

示例:
  urldb plugin create my_hook hook
  urldb plugin create add_user_table migration`,
	Args: cobra.MinimumNArgs(2),
	Run:  runCreatePlugin,
}

// listCmd 列出插件命令
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有插件文件",
	Long: `显示系统中所有的钩子文件和迁移文件`,
	Run:  runListPlugins,
}

// validateCmd 验证插件命令
var validateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "验证插件文件",
	Long: `验证插件文件的格式和内容

示例:
  urldb plugin validate ./hooks/my_hook.plugin.js`,
	Args: cobra.ExactArgs(1),
	Run:  runValidatePlugin,
}

// statsCmd 统计命令
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "显示插件统计信息",
	Long: `显示插件系统的统计信息，包括文件数量、大小等`,
	Run:  runPluginStats,
}

// InitPluginCommands 初始化插件命令
func InitPluginCommands() {
	pluginCmd.AddCommand(createCmd)
	pluginCmd.AddCommand(listCmd)
	pluginCmd.AddCommand(validateCmd)
	pluginCmd.AddCommand(statsCmd)
}

// runCreatePlugin 运行创建插件命令
func runCreatePlugin(cmd *cobra.Command, args []string) {
	pluginName := args[0]
	pluginType := args[1]

	pluginUtils := NewPluginUtils()

	if err := pluginUtils.CreatePluginTemplate(pluginName, pluginType); err != nil {
		utils.Error("创建插件模板失败: %v", err)
		os.Exit(1)
	}

	utils.Info("插件模板创建成功！")
	utils.Info("文件位置: %s", filepath.Join(
		map[string]string{
			"hook":      "./hooks",
			"migration": "./migrations",
		}[pluginType],
		pluginName+map[string]string{
			"hook":      ".plugin.js",
			"migration": ".js",
		}[pluginType],
	))
}

// runListPlugins 运行列出插件命令
func runListPlugins(cmd *cobra.Command, args []string) {
	pluginUtils := NewPluginUtils()

	if err := pluginUtils.ListPlugins(); err != nil {
		utils.Error("列出插件失败: %v", err)
		os.Exit(1)
	}
}

// runValidatePlugin 运行验证插件命令
func runValidatePlugin(cmd *cobra.Command, args []string) {
	filePath := args[0]

	pluginUtils := NewPluginUtils()

	if err := pluginUtils.ValidatePlugin(filePath); err != nil {
		utils.Error("插件验证失败: %v", err)
		os.Exit(1)
	}

	utils.Info("插件验证通过: %s", filePath)
}

// runPluginStats 运行插件统计命令
func runPluginStats(cmd *cobra.Command, args []string) {
	pluginUtils := NewPluginUtils()

	stats := pluginUtils.GetPluginStats()

	fmt.Println("=== 插件系统统计 ===")
	fmt.Printf("钩子文件数量: %d\n", stats["hooks_count"])
	fmt.Printf("迁移文件数量: %d\n", stats["migrations_count"])
	fmt.Printf("类型文件存在: %v\n", stats["types_file_exists"])
	fmt.Printf("最后更新时间: %s\n", stats["last_updated"])
}