//go:build linux && !windows

package contitle

// Set Windows console Title text (Header)
//
// base take here https://github.com/lxi1400/GoTitle
func SetTitle(title string) error {
	return nil
}
