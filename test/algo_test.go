package spritepacker

import (
	"fmt"
	"math/rand"
	"os"
	"spritepacker/pack"
	"strings"
	"testing"
	"time"
)

type AlgoResult struct {
	Rects    []pack.PackedRect
	Title    string
	FillRate float64
	TimeUsed int64
}

func Test_packRects(t *testing.T) {
	reqRects := generateRandomRects(100, 200, 200)
	options, err := pack.NewOptions().MaxSize(1024, 1024).AllowRotate(true).Validate()
	if err != nil {
		fmt.Println(err)
		return
	}
	var results []AlgoResult
	algorithms := []struct {
		name string
		algo pack.Algorithm
	}{
		{"Basic", pack.AlgoBasic},
		{"Skyline", pack.AlgoSkyline},
		{"MaxRects", pack.AlgoMaxRects},
	}
	for _, a := range algorithms {
		start := time.Now()
		packed := pack.NewPacker(options.Algorithm(a.algo)).PackRect(reqRects)
		elapsed := time.Since(start).Nanoseconds()
		fmt.Printf("%s FillRate: %.2f%%, Time: %d ns\n", a.name, packed.Bin.FillRate*100, elapsed)
		results = append(results, AlgoResult{
			Rects:    packed.Bin.PackedRects,
			Title:    a.name,
			FillRate: packed.Bin.FillRate,
			TimeUsed: elapsed,
		})
	}
	generateComparisonHTML(results)
}

func generateRandomRects(count, maxW, maxH int) []pack.Rect {
	rects := make([]pack.Rect, 0, count)
	for i := 0; i < count; i++ {
		if rand.Intn(2) == 0 {
			w := rand.Intn(maxW/2) + maxW/2
			h := rand.Intn(maxH/2) + maxH/2
			rects = append(rects, pack.NewRectById(w, h, i))
		} else {
			w := rand.Intn(maxW/3) + 1
			h := rand.Intn(maxH/3) + 1
			rects = append(rects, pack.NewRectById(w, h, i))
		}
	}
	return rects
}

func generateComparisonHTML(results []AlgoResult) {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Algorithm Comparison</title>
  <style>
    body { font-family: sans-serif; margin: 0; padding: 0; }
    h2 { text-align: center; }
    .canvas-container {
      display: flex;
      justify-content: space-around;
      align-items: center;
      margin-top: 20px;
      flex-wrap: wrap;
    }
    .panel {
      margin: 20px;
      text-align: center;
    }
    canvas { 
      border: 1px solid #ccc; 
      background: #fff;
      width: 100%;
      height: auto;
    }
    .info { margin-top: 10px; }
  </style>
</head>
<body>
  <h2>Algorithm Comparison</h2>
  <div class="canvas-container">
`
	for i, res := range results {
		canvasID := fmt.Sprintf("canvas%d", i)
		html += fmt.Sprintf(`
    <div class="panel">
      <h3>%s</h3>
      <canvas id="%s" width="400" height="400"></canvas>
      <div class="info">
        <p>FillRate: %.2f%%</p>
        <p>Time: %d ns</p>
      </div>
    </div>
`, res.Title, canvasID, res.FillRate*100, res.TimeUsed)
	}

	html += `</div><script>`
	for i, res := range results {
		canvasID := fmt.Sprintf("canvas%d", i)
		html += fmt.Sprintf(`
(function() {
  const data = %s;
  const isRotate = %s;
  const canvas = document.getElementById("%s");
  const ctx = canvas.getContext("2d");
  let maxX = 0, maxY = 0;
  data.forEach(rect => {
    const x2 = rect.x + rect.w;
    const y2 = rect.y + rect.l;
    if (x2 > maxX) maxX = x2;
    if (y2 > maxY) maxY = y2;
  });
  const scale = Math.min(canvas.width / maxX, canvas.height / maxY);
  data.forEach((rect, i) => {
    const color = "#" + Math.floor(Math.random()*16777215).toString(16).padStart(6, "0");
    const x = rect.x * scale;
    const y = rect.y * scale;
    const w = rect.w * scale;
    const h = rect.l * scale;
    ctx.fillStyle = color;
    ctx.fillRect(x, y, w, h);
    ctx.strokeStyle = "black";
    ctx.strokeRect(x, y, w, h);
    ctx.fillStyle = "black";
    ctx.font = "12px Arial";
    ctx.fillText(i + (isRotate[i] === 1 ? " (R)" : ""), x + 3, y + 12);
  });
})();
`, formatRects(res.Rects), formatRotateFlags(res.Rects), canvasID)
	}
	html += `</script></body></html>`
	_ = os.MkdirAll("output", 0755)
	err := os.WriteFile("output/compare.html", []byte(html), 0644)
	if err != nil {
		fmt.Println("write file error:", err)
	} else {
		fmt.Println("Generated " + "output/compare.html")
	}
}

func formatRects(rects []pack.PackedRect) string {
	parts := make([]string, len(rects))
	for i, r := range rects {
		parts[i] = fmt.Sprintf("{x:%d,y:%d,w:%d,l:%d}", r.X, r.Y, r.W, r.H)
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func formatRotateFlags(rects []pack.PackedRect) string {
	flags := make([]string, len(rects))
	for i, r := range rects {
		if r.IsRotated {
			flags[i] = "1"
		} else {
			flags[i] = "0"
		}
	}
	return "[" + strings.Join(flags, ",") + "]"
}
