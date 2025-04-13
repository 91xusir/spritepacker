------

# 🧩 SpritePacker 

**SpritePacker** is a 2D sprite sheet packing tool, designed for game developers 

------

## 📦 安装

```bash
go install github.com/91xusir/spritepacker@latest
```

------

## 🛠 使用方式

### CLI 命令行

```bash
go install github.com/91xusir/spritepacker@latest
spritepacker help
```

- ## 🛠️ 打包选项

  | 参数         | 说明                                                         |
  | ------------ | ------------------------------------------------------------ |
  | `-input`     | 输入目录，包含待打包的精灵图像（必需）                       |
  | `-output`    | 输出目录，默认为 `output`                                    |
  | `-maxw`      | 图集最大宽度，默认 `2048`                                    |
  | `-maxh`      | 图集最大高度，默认 `2048`                                    |
  | `-padding`   | 精灵之间的间距，默认 `2`                                     |
  | `-rotate`    | 是否允许旋转精灵以节省空间                                   |
  | `-pot`       | 强制图集尺寸为 2 的幂（Power Of Two）                        |
  | `-pma`       | 使用预乘透明度（Premultiplied Alpha）                        |
  | `-name`      | 输出图集的基础名称，默认 `atlas`                             |
  | `-sort`      | 打包前是否排序图像（默认启用）                               |
  | `-trim`      | 是否裁剪透明边缘（默认启用）                                 |
  | `-tolerance` | 裁剪透明边缘的容差（0~255，默认 `1`）                        |
  | `-same`      | 是否启用重复图像检测（默认启用）                             |
  | `-algo`      | 打包算法选择： 0 = Basic，1 = Skyline，2 = MaxRects          |
  | `-heuristic` | MaxRects 启发策略（仅用于 `-algo 2`）： 0 = BestShortSideFit 1 = BestLongSideFit 2 = BestAreaFit 3 = BottomLeft 4 = ContactPoint |

  ------

  ## 🛠️ 解包选项

  | 参数        | 说明                      |
  |-----------|-------------------------|
  | `-unpack` | 图集对应的 JSON 路径（必需）       |
  | `-img`    | 图集图像路径（可选，默认按 JSON 名推导） |
  | `-output` | 解包输出目录（可选，默认按 JSON 名推导） |

  ------

  ## 📦 示例

  ```bash
  使用默认配置打包 ./sprites 下的所有图像
  spritepacker -input ./sprites -output ./atlases
  
  # 自定义图集尺寸、间距、开启旋转
  spritepacker -input ./sprites -maxw 1024 -maxh 1024 -padding 4 -rotate
  
  # 解包 atlas.json 对应图集
  spritepacker -unpack ./atlases/atlas.json -output ./unpacked
  ```

------

## 🧪 API

```bash
go get -u github.com/91xusir/spritepacker
```

使用示例

```go
package main
import (
	"github.com/91xusir/spritepacker/pack"
	"encoding/json"
	"os"
	"path/filepath"
)
func Pack() {
	// Step 1: Create packing options and set the maximum atlas size to 4096x4096
	options := pack.NewOptions().MaxSize(4096, 4096)

	// Optional: Validate the options. If invalid, an error will be returned.
	// Here we ignore the error for demonstration purposes.
	options, _ = options.Validate()

	// Step 2: Create the sprite packer using the configured options
	spritePacker := pack.NewPacker(options)

	// Step 3: Collect all sprite image paths from the "input" directory
	spriteImgPaths, _ := pack.GetFilesInDirectory("input")

	// Step 4: Pack all sprite images and generate atlas metadata and images
	spriteAtlasInfo, atlasImages, _ := spritePacker.PackSprites(spriteImgPaths)
	
	// Step 5: Save each atlas image to the "output" directory
	for i := range atlasImages {
		// If you plan to unpack the atlas using spriteAtlasInfo later,
		// you should use the atlas name from spriteAtlasInfo as the filename.
		// Only the name is required, and the file extension will be determined
		// based on the format passed (e.g., pack.PNG, pack.JPEG, etc.).
		// WithCLV is function to set the compression level of the image.
		outputPath := filepath.Join("output", spriteAtlasInfo.Atlases[i].Name)
		_ = pack.SaveImg(outputPath, atlasImages[i], pack.PNG, pack.WithCLV(pack.DefaultCompression))
	}

	// Step 6: Save the atlas metadata as a JSON file
	// By default, all atlas information is saved into a single file.
	// You can read spriteAtlasInfo to save each atlas information separately.
	// Correspondingly, you would need to customize the UnpackSprites method
	// to read each atlas information and save each atlas image.
	jsonBytes, _ := json.MarshalIndent(spriteAtlasInfo, "", "  ")
	_ = os.WriteFile("output/atlas.json", jsonBytes, os.ModePerm)

}

func Unpack() {
	// default output and atlas images path is the same as the atlas.json path
	// if you want to change this, you can use the following code:
	// pack.UnpackSprites("output/atlas.json", pack.WithAtlasImgPath("output"), pack.WithOutputPath("output"))
	err := pack.UnpackSprites("output/atlas.json")
	if err != nil {
		return
	}
	// Then you can use TestImageDiff to compare the output images with the input images
	// to verify that the unpacking process is correct.
}
```

