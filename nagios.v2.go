//V2. 30 July 2024 8:56 AM
//Herwin Yudha Setyawan

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/go-gomail/gomail"
	"github.com/sirupsen/logrus"
)

func main() {
	// Original URL
	rawURL := "nagios-link"

	// Parse the URL
	u, err := url.Parse(rawURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	// Get username and password
	username := "secret"
	password := "secret"

	// Reconstruct the URL with properly encoded credentials
	u.User = url.UserPassword(username, password)
	encodedURL := u.String()

	emailTo := []string{"List Email To"}
	emailCC := []string{"List Email CC"}

	// Prompt user to choose option
	fmt.Println("")
	fmt.Println("<><><><><><><><><><><><><><><><><><><><><><><><><><><><>")
	fmt.Println("Automation Nagios Mail Report")
	fmt.Println("Options:")
	fmt.Println("1. Run script scheduled time at 15:00 WIB or 10:00 CEST")
	fmt.Println("2. Run script now at current time")
	fmt.Println("<><><><><><><><><><><><><><><><><><><><><><><><><><><><>")
	fmt.Print("Choose option: ")
	var choice int
	fmt.Scan(&choice)

	switch choice {
	case 1:
		for {
			// Wait until the target time (15:00 WIB / 10:00 CEST) "Asia/Jakarta"
			err := runAtTargetTime(encodedURL, emailTo, emailCC)
			if err != nil {
				log.Printf("Failed to run script at target time: %v", err)
				log.Println("Retrying in 2 minutes...")
				log.Println("Please check what is the error, check vpn 3uk....")
				time.Sleep(2 * time.Minute)
				err := runNow(encodedURL, emailTo, emailCC)
				if err != nil {
					log.Fatalf("Failed to run script now: %v", err)
					log.Println("Please check what is the error....")
				}
				continue
			}
			break
		}
	case 2:
		err := runNow(encodedURL, emailTo, emailCC)
		if err != nil {
			log.Fatalf("Failed to run script now: %v", err)
			log.Println("Please check what is the error....")
		}
		// Return to default schedule
		for {
			err = runAtTargetTime(encodedURL, emailTo, emailCC)
			if err != nil {
				log.Printf("Failed to run script at target time after running now: %v", err)
				log.Println("Retrying in 2 minutes...")
				log.Println("Please check what is the error, check vpn 3uk....")
				time.Sleep(2 * time.Minute)
				continue
			}
			break
		}
	default:
		fmt.Println("Invalid choice")
	}
}

func runAtTargetTime(encodedURL string, emailTo, emailCC []string) error {
	waitUntilTargetTime(15, 0, "Asia/Jakarta") // Use appropriate timezone

	// Set up logging
	logger := logrus.New()
	file, err := os.OpenFile("Error-log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()
	logger.Out = file

	// Take screenshot
	log.Println("Taking screenshot...")
	screenshot, err := takeScreenshot(encodedURL, logger)
	if err != nil {
		logger.Errorf("Failed to take screenshot: %v", err)
		log.Println("Error: Failed to take screenshot.")
		return err
	}

	// Send email with screenshot
	log.Println("Sending email...")
	err = sendEmail(emailTo, emailCC, screenshot)
	if err != nil {
		logger.Errorf("Failed to send email: %v", err)
		log.Println("Error: Failed to send email.")
		return err
	}

	// Log success
	successLog, err := os.OpenFile("Success-log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Errorf("Failed to open success log file: %v", err)
		log.Println("Error: Failed to open success log file.")
		return err
	}
	defer successLog.Close()
	successLogger := log.New(successLog, "", log.LstdFlags)
	successLogger.Printf("Email sent successfully at %v", time.Now())
	log.Println("Email sent successfully.")

	return nil
}

func runNow(encodedURL string, emailTo, emailCC []string) error {
	log.Println("Running script now...")

	// Set up logging
	logger := logrus.New()
	file, err := os.OpenFile("Error-log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()
	logger.Out = file

	// Take screenshot
	log.Println("Taking screenshot...")
	screenshot, err := takeScreenshot(encodedURL, logger)
	if err != nil {
		logger.Errorf("Failed to take screenshot: %v", err)
		log.Println("Error: Failed to take screenshot.")
		return err
	}

	// Send email with screenshot
	log.Println("Sending email...")
	err = sendEmail(emailTo, emailCC, screenshot)
	if err != nil {
		logger.Errorf("Failed to send email: %v", err)
		log.Println("Error: Failed to send email.")
		return err
	}

	// Log success
	successLog, err := os.OpenFile("Success-log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Errorf("Failed to open success log file: %v", err)
		log.Println("Error: Failed to open success log file.")
		return err
	}
	defer successLog.Close()
	successLogger := log.New(successLog, "", log.LstdFlags)
	successLogger.Printf("Email sent successfully at %v", time.Now())
	log.Println("Email sent successfully.")

	return nil
}

func waitUntilTargetTime(targetHour, targetMinute int, timezone string) {
	for {
		now := time.Now()
		location, _ := time.LoadLocation(timezone)
		targetTime := time.Date(now.Year(), now.Month(), now.Day(), targetHour, targetMinute, 0, 0, location)

		// If the target time has already passed today, set it for tomorrow
		if now.After(targetTime) {
			targetTime = targetTime.Add(24 * time.Hour)
		}

		// Calculate the duration to wait
		duration := targetTime.Sub(now)
		fmt.Printf(" \n")
		fmt.Printf("<><><>h<>y<>s<>s<>p<>y<><><><><><><><><><><><>\n")
		fmt.Printf("<><>.....Please Don't Close this..........<><>\n")
		fmt.Printf("<><>......Automate Nagios Mail Report.....<><>\n")
		fmt.Printf("<><>..Please..check.if.VPN...Connected....<><>\n")
		fmt.Printf("<><>....File..in..Desktop...nagios.exe....<><>\n")
		fmt.Printf("<><><><><><><><><><><><><><><><><><><><><><><>\n")
		fmt.Printf(" \n")
		fmt.Printf("Waiting until %s (%v from now)\n", targetTime, duration)

		// Wait until the target time
		time.Sleep(duration)

		// Check if the current time matches the target time to avoid slight delays
		if time.Now().After(targetTime) {
			break
		}
	}
}

func takeScreenshot(encodedURL string, logger *logrus.Logger) ([]byte, error) {
	var screenshot []byte
	retries := 3
	for i := 0; i < retries; i++ {
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		err := chromedp.Run(ctx, fullScreenshot(encodedURL, 90, &screenshot))
		if err == nil {
			return screenshot, nil
		}

		logger.Errorf("Attempt %d: failed to take screenshot: %v", i+1, err)
		log.Printf("Attempt %d: failed to take screenshot: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to take screenshot after %d attempts", retries)
}

func fullScreenshot(encodedURL string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(encodedURL),
		chromedp.Sleep(10 * time.Second), // Increase the sleep time to allow the page to load completely
		chromedp.FullScreenshot(res, quality),
	}
}

func sendEmail(to, cc []string, attachment []byte) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "Email Gateway")
	m.SetHeader("To", to...)
	if len(cc) > 0 {
		m.SetHeader("Cc", cc...)
	}
	currentTime := time.Now().Format("02/01/2006")
	subject := fmt.Sprintf("Nagios Daily Report %s", currentTime)
	m.SetHeader("Subject", subject)

	body := `Hello all,<br><br>

Please find below the current status of the platform.<br>
Take into account that the severity reported is internal; not customer facing.<br><br>

Nagios current status:<br>
<img src="cid:nagios.png">
`
	m.SetBody("text/html", body)

	// Save attachment to file to check if the screenshot is valid
	err := os.WriteFile("nagios_check.png", attachment, 0644) // Check the screenshot file
	if err != nil {
		return fmt.Errorf("failed to write attachment to file: %v", err)
	}

	// Attach screenshot as inline image
	m.Attach("nagios.png", gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(attachment)
		if err != nil {
			return fmt.Errorf("failed to write attachment to email: %v", err)
		}
		return nil
	}), gomail.SetHeader(map[string][]string{
		"Content-ID": {"<nagios.png>"},
	}))

	d := gomail.NewDialer("smtp", smtp-port, "secretmail@mail.com", "secretpassword")

	// Log email details for debugging
	log.Printf("Sending email to: %s, cc: %s, with subject: %s", to, cc, subject)

	err = d.DialAndSend(m)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
