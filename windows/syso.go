package windows

import (
	"fmt"
	"os"

	"github.com/hallazzang/syso"
	"github.com/hallazzang/syso/pkg/coff"
	"github.com/intex-software/goutils/internal"
)

type SysOConfig struct {
	Language string
	Charset  string

	CompanyName      string
	FileDescription  string
	InternalName     string
	LegalCopyright   string
	OriginalFilename string
	ProductName      string
	FileVersion      string
	ProductVersion   string

	Manifest string
	Icon     string
}

func CreateSysO(outFile string, config SysOConfig) (err error) {
	c := coff.New()
	if err := syso.EmbedIcon(c, &syso.FileResource{ID: 1, Path: config.Icon}); err != nil {
		return fmt.Errorf("failed to embed icon #%d: %v", 1, err)
	}
	if err := syso.EmbedManifest(c, &syso.FileResource{ID: 2, Path: config.Manifest}); err != nil {
		return fmt.Errorf("failed to embed manifest %v", err)
	}

	vi := syso.VersionInfoResource{
		ID: internal.Ptr(1),
		Fixed: &syso.VersionInfoFixed{
			FileVersion:    internal.Ptr(config.FileVersion),
			ProductVersion: internal.Ptr(config.ProductVersion),
		},
		StringTables: []*syso.VersionInfoStringTable{
			{

				Language: internal.Ptr(config.Language),
				Charset:  internal.Ptr(config.Charset),
				Strings: &syso.VersionInfoStrings{
					CompanyName:      internal.Ptr(config.CompanyName),
					FileDescription:  internal.Ptr(config.FileDescription),
					InternalName:     internal.Ptr(config.InternalName),
					LegalCopyright:   internal.Ptr(config.LegalCopyright),
					OriginalFilename: internal.Ptr(config.OriginalFilename),
					ProductName:      internal.Ptr(config.ProductName),
					FileVersion:      internal.Ptr(config.FileVersion),
					ProductVersion:   internal.Ptr(config.ProductVersion),
				},
			},
		},
		Translations: []*syso.VersionInfoTranslation{
			{
				Language: internal.Ptr(config.Language),
				Charset:  internal.Ptr(config.Charset),
			},
		},
	}
	if err := syso.EmbedVersionInfo(c, &vi); err != nil {
		return fmt.Errorf("failed to embed version info %v", err)
	}

	w, err := os.Create(outFile)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	if _, err := c.WriteTo(w); err != nil {
		panic(err)
	}

	return
}
