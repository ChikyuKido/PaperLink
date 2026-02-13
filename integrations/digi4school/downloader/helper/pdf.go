package helper

import (
	"fmt"
	"github.com/signintech/gopdf"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ConvertAndCompressSVG(downloadPath, filename string) (string, error) {
	inputSVG := filepath.Join(downloadPath, filename)
	outputPDF := filepath.Join(downloadPath, strings.TrimSuffix(filename, ".svg")+".pdf")
	tmpPDF := outputPDF + ".tmp"
	finalPDF := outputPDF + ".final"

	if err := exec.Command(
		"rsvg-convert",
		"-f", "pdf",
		"-o", outputPDF,
		inputSVG,
	).Run(); err != nil {
		return "", fmt.Errorf("rsvg-convert failed: %w", err)
	}

	if err := exec.Command(
		"gs",
		"-sDEVICE=pdfwrite",
		"-dCompatibilityLevel=1.5",
		"-dDetectDuplicateImages=true",
		"-dSubsetFonts=true",
		"-dEmbedAllFonts=true",
		"-dCompressFonts=true",
		"-dNOPAUSE",
		"-dQUIET",
		"-dBATCH",
		"-sOutputFile="+tmpPDF,
		outputPDF,
	).Run(); err != nil {
		return "", fmt.Errorf("ghostscript failed: %w", err)
	}

	if err := exec.Command(
		"qpdf",
		"--object-streams=generate",
		"--stream-data=compress",
		tmpPDF,
		finalPDF,
	).Run(); err != nil {
		return "", fmt.Errorf("qpdf failed: %w", err)
	}

	if err := os.Rename(finalPDF, outputPDF); err != nil {
		return "", err
	}
	_ = os.Remove(tmpPDF)

	return outputPDF, nil
}

func ConvertPNGtoPDF(downloadPath, filename string, pngWidthPx, pngHeightPx int) (string, error) {
	inputPNG := filepath.Join(downloadPath, filename)
	outputPDF := filepath.Join(downloadPath, strings.TrimSuffix(filename, ".png")+".pdf")
	pdfWidthPt, pdfHeightPt := 595.0, 842.0
	dpi := 96.0
	imgW := float64(pngWidthPx) * 72.0 / dpi
	imgH := float64(pngHeightPx) * 72.0 / dpi

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: pdfWidthPt, H: pdfHeightPt}})
	pdf.AddPage()

	scaleW := pdfWidthPt / imgW
	scaleH := pdfHeightPt / imgH
	scale := scaleW
	if scaleH < scaleW {
		scale = scaleH
	}
	drawW := imgW * scale
	drawH := imgH * scale
	offsetX := (pdfWidthPt - drawW) / 2
	offsetY := (pdfHeightPt - drawH) / 2

	pdf.Image(inputPNG, offsetX, offsetY, &gopdf.Rect{W: drawW, H: drawH})

	if err := pdf.WritePdf(outputPDF); err != nil {
		return "", err
	}

	return outputPDF, nil
}
