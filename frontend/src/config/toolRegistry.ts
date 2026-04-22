export type ToolKey =
  | 'text-dedupe'
  | 'random-string'
  | 'text-diff'
  | 'json-tool'
  | 'markdown-editor'
  | 'mermaid-editor'
  | 'timestamp'
  | 'bookmarks'
  | 'secret-manager'
  | 'devtools'
  | 'qr-code'
  | 'qr-decode'
  | 'mind-map'
export type ToolCategory = '文本工具' | '开发工具' | '图像工具' | '图形工具'

export interface ToolMenuItem {
  key: ToolKey
  title: string
  description: string
  category: ToolCategory
}

const baseTools: ToolMenuItem[] = [
  {
    key: 'text-dedupe',
    title: '文本去重',
    description: '按行删除重复内容',
    category: '文本工具',
  },
  {
    key: 'random-string',
    title: '随机字符串',
    description: '生成密码和 Token',
    category: '文本工具',
  },
  {
    key: 'text-diff',
    title: '文本差异对比',
    description: '查看新增和删除行',
    category: '文本工具',
  },
  {
    key: 'json-tool',
    title: 'JSON 工具',
    description: '格式化压缩和校验',
    category: '文本工具',
  },
  {
    key: 'markdown-editor',
    title: 'Markdown 编辑器',
    description: '实时编辑和预览',
    category: '文本工具',
  },
  {
    key: 'mermaid-editor',
    title: 'Mermaid 编辑器',
    description: '图表源码实时预览',
    category: '图形工具',
  },
  {
    key: 'timestamp',
    title: '时间戳转换',
    description: '时间和时间戳互转',
    category: '开发工具',
  },
  {
    key: 'bookmarks',
    title: '浏览器书签',
    description: '保存和打开常用网站',
    category: '开发工具',
  },
  {
    key: 'secret-manager',
    title: '账号密钥管理',
    description: '本地加密保存密码和 Key',
    category: '开发工具',
  },
  ...(import.meta.env.DEV
    ? [
      {
        key: 'devtools' as const,
        title: '开发者调试',
        description: '打开 DevTools 和查看环境',
        category: '开发工具' as const,
      },
    ]
    : []),
  {
    key: 'qr-code',
    title: '生成二维码',
    description: '文本和链接转 PNG',
    category: '图像工具',
  },
  {
    key: 'qr-decode',
    title: '解析二维码',
    description: '图片提取二维码内容',
    category: '图像工具',
  },
  {
    key: 'mind-map',
    title: '脑图',
    description: '缩进文本生成脑图',
    category: '图形工具',
  },
]

export const toolRegistry: ToolMenuItem[] = baseTools

export function groupTools(tools: ToolMenuItem[]) {
  return tools.reduce<Array<{ category: ToolCategory; tools: ToolMenuItem[] }>>((groups, tool) => {
    const group = groups.find((item) => item.category === tool.category)
    if (group) {
      group.tools.push(tool)
    } else {
      groups.push({ category: tool.category, tools: [tool] })
    }
    return groups
  }, [])
}
