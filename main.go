package main

import (
    "context"
    //"io/ioutil"
    "log"
    "net/http"
    "time"

    "github.com/chromedp/cdproto/emulation"
    "github.com/chromedp/cdproto/page"
    "github.com/chromedp/chromedp"
)

func main() {
    // Create a new HTTP server
    http.HandleFunc("/generate-pdf", generatePDFHandler)

    // Start the server
    log.Println("Server started on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func generatePDFHandler(w http.ResponseWriter, r *http.Request) {
    // Obtain the URL from the query parameter
    url := r.URL.Query().Get("url")
    if url == "" {
        http.Error(w, "Missing 'url' query parameter", http.StatusBadRequest)
        return
    }

    // Create a new chromedp context
    ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()

    // Chained chromedp tasks to generate the PDF
    var pdfBuffer []byte
    err := chromedp.Run(ctx, generatePDF(url, &pdfBuffer))
    if err != nil {
        http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
        log.Println(err)
        return
    }

    // Set the response headers
    w.Header().Set("Content-Type", "application/pdf")
    //w.Header().Set("Content-Disposition", "attachment; filename=output.pdf")
    //w.Header().Set("Content-Length", strconv.Itoa(len(pdfBuffer)))

    // Write the PDF buffer to the response
    if _, err := w.Write(pdfBuffer); err != nil {
        http.Error(w, "Failed to write response", http.StatusInternalServerError)
        log.Println(err)
        return
    }
}

func generatePDF(url string, pdfBuffer *[]byte) chromedp.Tasks {
    return chromedp.Tasks{
        emulation.SetDeviceMetricsOverride(1920, 1080, 1.0, false),
        chromedp.Navigate(url),
        chromedp.Sleep(2 * time.Second),
        chromedp.ActionFunc(func(ctx context.Context) error {
            buf, _, err := page.PrintToPDF().Do(ctx)
            if err != nil {
                return err
            }
            *pdfBuffer = buf
            return nil
        }),
    }
}