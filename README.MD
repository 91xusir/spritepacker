------

# 🧩 SpritePacker | 精灵图打包工具

**SpritePacker** is a 2D sprite sheet packing tool, designed for game developers 

------

## 📦 Installation | 安装

```
go install github.com/yourname/spritepacker@latest
```

or

```
git clone https://github.com/yourname/spritepacker.git
cd spritepacker
go build .
```

------

## 🛠 Usage | 使用方式

### CLI 命令行

```bash
spritepacker -input ./sprites -output ./atlas -maxw 1024 -maxh 1024 -trim -rotate
```

### Options 参数说明

| Flag        | Description (EN)                  | 描述（中文）           |
|-------------|-----------------------------------|------------------|
| `-input`    | Input image folder                | 输入图片文件夹          |
| `-output`   | Output atlas & metadata directory | 输出图集和元数据目录       |
| `-maxw`     | Maximum atlas width               | 最大图集宽度           |
| `-maxh`     | Maximum atlas height              | 最大图集高度           |
| `-trim`     | Trim transparent edges            | 裁剪透明边缘           |
| `-rotate`   | Allow 90-degree rotation          | 允许 90 度旋转优化      |
| `-autosize` | Auto resize the canvas            | 自动适配画布尺寸         |
| `-format`   | Output format (json, xml, etc.)   | 输出格式（json、xml 等） |

------

## 📂 Output Example | 输出示例

```
├── atlas.png        # Packed texture 图集
├── atlas.json       # Metadata 元数据（包含坐标、尺寸等）
```

------

## 🧠 Algorithm | 算法

- **Basic**
- **MaxRects**
- **Skyline**

------

## 🧪 Go API Example | Go 语言调用示例

```GO
func Test_Packing(t *testing.T) {
	// you can use func NewOptions to set options
	options := packer.NewOptions(
		packer.WithMaxSize(1024, 1024),
		packer.WithPadding(0),
		packer.WithAlgorithm(packer.AlgoSkyline),
	)
	_ = options
	// you also can use OptionBuilder to set options
	options = packer.NewOptionBuilder().
		MaxSize(1024, 1024).
		Padding(0).
		Algorithm(packer.AlgoSkyline).
		Build()

	reqPackedRect, _ := getTestData("test_rect/rect.txt")

	packer := packer.NewPacker(options).SetRects(reqPackedRect)

	packResult := packer.Packing()

	fmt.Println(packResult)
}
```

