package progress_bar

import (
	"fmt"
	"io"
	"time"
)

// ProgressWriter tracks bytes written and throttles display updates to avoid flickering
type ProgressWriter struct {
	Total      int64
	Current    int64
	Label      string
	lastUpdate time.Time
	interval   time.Duration
}

// NewProgressWriter creates a new ProgressWriter with a custom update interval (e.g., 200ms - 1s)
func NewProgressWriter(total int64, label string, interval time.Duration) *ProgressWriter {
	if interval <= 0 {
		interval = 500 * time.Millisecond
	}
	pw := &ProgressWriter{
		Total:    total,
		Label:    label,
		interval: interval,
	}
	pw.render(0, false)
	return pw
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.Current += int64(n)

	now := time.Now()
	// Update every interval or when completed
	if now.Sub(pw.lastUpdate) >= pw.interval || (pw.Total > 0 && pw.Current >= pw.Total) {
		pw.lastUpdate = now
		pw.render(pw.Current, pw.Total > 0 && pw.Current >= pw.Total)
	}
	return n, nil
}

func (pw *ProgressWriter) render(current int64, done bool) {
	if pw.Total <= 0 {
		fmt.Printf("\r%s: %s transferred...", pw.Label, formatBytes(current))
		if done {
			fmt.Println()
		}
		return
	}

	percent := float64(current) / float64(pw.Total) * 100
	if percent > 100 {
		percent = 100
	}

	// \r returns cursor to the beginning of the line to overwrite previous output
	fmt.Printf("\r%s: [%.1f%%] %s / %s", pw.Label, percent, formatBytes(current), formatBytes(pw.Total))
	if done {
		fmt.Println() // Print newline when 100% complete
	}
}

// Finish ensures the progress display shows 100% and ends with a newline
func (pw *ProgressWriter) Finish() {
	pw.render(pw.Current, true)
}

// TeeWriter creates a writer that duplicates its writes to w and the progress tracker
func (pw *ProgressWriter) Tee(w io.Writer) io.Writer {
	return io.MultiWriter(w, pw)
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
