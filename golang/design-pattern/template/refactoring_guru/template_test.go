package refactoring_guru

import "fmt"

func Example_Decorator() {
	// otp := otp{}

	// smsOTP := &sms{
	//  otp: otp,
	// }

	// smsOTP.genAndSendOTP(smsOTP, 4)

	// emailOTP := &email{
	//  otp: otp,
	// }
	// emailOTP.genAndSendOTP(emailOTP, 4)
	// fmt.Scanln()

	smsOTP := &sms{}
	o := otp{
		iOtp: smsOTP,
	}
	o.genAndSendOTP(4)

	fmt.Println("")
	emailOTP := &email{}
	o = otp{
		iOtp: emailOTP,
	}

	//Output:
	//SMS: generating random otp 1234
	//SMS: saving otp: 1234 to cache
	//SMS: sending sms: SMS OTP for login is 1234
	//SMS: publishing metrics
}
