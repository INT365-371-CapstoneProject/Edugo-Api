package utils

import (
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"time"
)

// GenerateOTP สร้างรหัส OTP แบบสุ่ม 6 หลัก
func GenerateOTP() string {
	// สร้าง random source แยกต่างหาก
	randomSource := rand.NewSource(time.Now().UnixNano())
	random := rand.New(randomSource)
	
	digits := "0123456789"
	otp := make([]byte, 6)
	for i := range otp {
		otp[i] = digits[random.Intn(len(digits))]
	}
	return string(otp)
}

// SendOTPEmail ส่ง OTP ไปยังอีเมล์
func SendOTPEmail(toEmail, otp string) error {
	from := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")
	to := []string{toEmail}
	
	// สร้าง HTML template สำหรับ email
	htmlBody := `
<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333333;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f9f9f9;
            border-radius: 5px;
        }
        .header {
            text-align: center;
            padding: 20px;
            background-color: #355FFF;
            color: white;
            border-radius: 5px 5px 0 0;
        }
        .content {
            padding: 20px;
            background-color: white;
            border-radius: 0 0 5px 5px;
        }
        .otp-code {
            font-size: 32px;
            font-weight: bold;
            text-align: center;
            color: #355FFF;
            padding: 10px;
            margin: 20px 0;
            letter-spacing: 5px;
        }
        .footer {
            text-align: center;
            font-size: 12px;
            color: #666666;
            margin-top: 20px;
        }
        .warning {
            color: #ff0000;
            font-size: 14px;
            margin-top: 15px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>รหัส OTP สำหรับรีเซ็ตรหัสผ่าน</h1>
        </div>
        <div class="content">
            <p>เรียนผู้ใช้งาน,</p>
            <p>คุณได้ขอรหัส OTP สำหรับการรีเซ็ตรหัสผ่าน กรุณาใช้รหัสด้านล่างนี้:</p>
            
            <div class="otp-code">
                %s
            </div>
            
            <p>รหัส OTP นี้จะหมดอายุใน 15 นาที</p>
            
            <div class="warning">
                ⚠️ โปรดอย่าเปิดเผยรหัส OTP นี้กับผู้อื่น
            </div>
        </div>
        <div class="footer">
            <p>อีเมลนี้ถูกส่งโดยอัตโนมัติ กรุณาอย่าตอบกลับ</p>
            <p>&copy; 2024 Edugo. All rights reserved.</p>
        </div>
    </div>
</body>
</html>`

	// สร้าง email message
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: รีเซ็ตรหัสผ่าน - Reset Password OTP\n"
	message := []byte(subject + mime + fmt.Sprintf(htmlBody, otp))

	// ส่ง email
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, message)
	
	if err != nil {
		log.Printf("Error details: %v", err)
		return err
	}
	
	return nil
}

// ValidateOTP ตรวจสอบว่า OTP ถูกต้องและยังไม่หมดอายุ
func ValidateOTP(storedOTP string, inputOTP string, expiredAt time.Time) bool {
	return storedOTP == inputOTP && time.Now().Before(expiredAt)
}
