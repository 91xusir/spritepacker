package spritepacker

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"spritepacker/pack"
	"strconv"
	"strings"
	"testing"
	"time"
)

type AlgoResult struct {
	Rects    []pack.PackedRect
	Title    string
	FillRate float64
	TimeUsed int64
	totalS   int
	usedS    int
	w        int
	h        int
}

// used fixed data and fixed size to test
func TestFixedDataFixedSize(t *testing.T) {
	reqRects := make([]pack.Rect, 100)
	for i := 0; i < 100; i++ {
		reqRects[i] = pack.Rect{W: 64, H: 64}
	}

	options := pack.NewOptions().
		MaxSize(512, 512).
		AllowRotate(true).
		Heuristic(pack.BestLongSideFit)

	results := packedWithAllAlgorithms(t, reqRects, options)

	generateComparisonHTML(results, "fixedData_fixedSize")
}

// used fixed data and different size to test
func TestFixedDataDiffSize(t *testing.T) {
	reqRects, _ := getTestData("test_rect.txt")

	options := pack.NewOptions().
		MaxSize(1024, 1024).AutoSize(true).
		AllowRotate(true)

	results := packedWithAllAlgorithms(t, reqRects, options)

	generateComparisonHTML(results, "fixedData_diffSize")
}

// used random data and random size to test
func TestRandomDateRandomSize(t *testing.T) {
	reqRects := generateRandomRects(100, 100, 100)

	options := pack.NewOptions().
		MaxSize(512, 512).
		AllowRotate(true).
		Heuristic(pack.BestLongSideFit)

	results := packedWithAllAlgorithms(t, reqRects, options)

	generateComparisonHTML(results, "randomDate_randomSize")
}

func packedWithAllAlgorithms(t *testing.T, reqRects []pack.Rect, options *pack.Options) []AlgoResult {
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
		r := slices.Clone(reqRects)
		packed := pack.NewPacker(options.Algorithm(a.algo)).PackRect(r)
		elapsed := time.Since(start).Nanoseconds()
		t.Logf("%s FillRate: %.2f%%, Time: %d ns\n", a.name, packed.Bin.FillRate*100, elapsed)
		results = append(results, AlgoResult{
			Rects:    packed.Bin.PackedRects,
			Title:    a.name,
			FillRate: packed.Bin.FillRate,
			TimeUsed: elapsed,
			totalS:   packed.Bin.Area(),
			usedS:    packed.Bin.UsedArea,
			w:        packed.Bin.W,
			h:        packed.Bin.H,
		})
	}
	return results
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

func generateComparisonHTML(results []AlgoResult, filename string) {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Algorithm Comparison</title>
  <style>
    body { font-family: sans-serif; margin: 0; padding: 0; zoom: 0.6}
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
      <canvas id="%s" width="%d" height="%d"></canvas>
      <div class="info">
		<p>UsedArea: %d</p>
		<p>TotalArea: %d</p>
        <p>FillRate: %.2f%%</p>
        <p>Time: %d ns</p>
		<p>Count: %d</p>
      </div>
    </div>
`, res.Title, canvasID, res.w, res.h, res.usedS, res.totalS, res.FillRate*100, res.TimeUsed, len(res.Rects))
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
  data.forEach((rect, i) => {
    const color = "#" + Math.floor(Math.random()*16777215).toString(16).padStart(6, "0");
    const x = rect.x;
    const y = rect.y;
    const w = rect.w;
    const h = rect.l;
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
	html += `
	   const panels = document.querySelectorAll(".panel");
	   const maxPageWidth = 1200;
	   let totalContentWidth = 0;
	   panels.forEach(panel => {
	     totalContentWidth += panel.offsetWidth;
	   });
	   const zoom = Math.min(1, maxPageWidth / totalContentWidth);
	   document.body.style.zoom = zoom.toString();
	`
	html += `</script></body></html>`
	_ = os.MkdirAll("output", 0755)
	path := filepath.Join("output", filename+".html")
	err := os.WriteFile(path, []byte(html), 0644)
	if err != nil {
		return
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

func getTestData(path string) ([]pack.Rect, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var reqPackedRect []pack.Rect
	id := 0
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		w, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("an error in parsing rectangle width: %w", err)
		}
		h, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("an error in parsing rectangle height: %w", err)
		}
		rect := pack.NewRectById(w, h, id)
		id++
		reqPackedRect = append(reqPackedRect, rect)

	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return reqPackedRect, nil
}
