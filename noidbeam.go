//DOES THE CARRIER_FREQ ACTUALLY NEED TO BE A MULTIPLE OF DATA_FREQ??? YOU SHOULD PROBABLY KNOW THAT
package main

import "wav"
import "fmt"


var DATA_FREQ = 100 //frequency that the data is sent in in hz
var CARRIER_FREQ = 4000 //frequency of carrier signal
var SAMPLE_RATE = 40000 //sample rate of the output, //MUST BE AN INTEGER MULTIPLE OF CARRIER_FREQ AND DATA_FREQ

var SILENCE = byte(128) //NOT 0 OR 127; its 128 because amplitude is -128 to 127, NOT -127 to 127


func toBytes(data []int) []byte {
	bytes := make([]byte, len(data))
	for e, b := range data {
		bytes[e] = byte(b + 128)
	}
	return bytes
}

func byte2Bits(b byte) []uint8 {
	mask := uint8(1)
	bits := make([]uint8, 8)
	for i := uint8(0); i < 8; i++ {
		if b & (mask << i) != 0 {
			bits[7 - i] = 1
		} //dont need to set any bits to 0 because make() initialized them to 0 by default
	}
	return bits
}

func bytesToBinary(bytes []byte) []byte {
	allBits := []byte{}
	currentBits := []byte{}
	for _, b := range bytes {
		currentBits = byte2Bits(b)
		currentBits = append([]byte{1}, currentBits...) //ADD START BIT
		currentBits = append(currentBits, 0) //ADD END BIT
		allBits = append(allBits, currentBits...)
	}
	return allBits
}

func carrier(length int) []byte {
	data := make([]byte, length)
	isHigh := false
	halfPeriod := SAMPLE_RATE / CARRIER_FREQ
	for i := 0; i < length; i++ {
		if i % halfPeriod == 0 { //every half a period
			isHigh = !isHigh //flip isHigh
		}
		if isHigh { //if its high, make the sample high
			data[i] = 255
		} else { //otherwise make it silence
			data[i] = SILENCE
		}
	}
	return data
}

func modulate(carrier []byte, data []byte, bitLength int) []byte {
	carrierLoc := 0
	for _, b := range data {
		if b == 0 {
			for i := 0; i < bitLength; i++ { //make the next bit length of samples all silent
				carrier[carrierLoc + i] = SILENCE
			}
		}
		carrierLoc += bitLength //move to the next data sample
	}
	return carrier
}

func addSilence(data []byte) []byte { //NEED TO ADD SILENCE BECAUSE WHEN I PLAY THE AUDIO IN WINDOWS MEDIA PLAYER IT FADES IN AND OUT
	silence := make([]byte, SAMPLE_RATE) //because sample rate is just the number of samples in a second
	for e := range silence {
		silence[e] = SILENCE
	}
	data = append(silence, data...)
	return append(data, silence...)
}

func checkGlobals() bool {
	if SAMPLE_RATE % CARRIER_FREQ != 0 || SAMPLE_RATE % DATA_FREQ != 0 || CARRIER_FREQ % DATA_FREQ != 0 {
		fmt.Println("CARRIER OR DATA FREQUENCY IS NOT MULTIPLE OF SAMPLE_RATE")
		return false
	}
	return true
}

/*
func getGlobals() {
	fmt.Print("Data Frequency: ")
	DATA_FREQ, _ = fmt.Scan()
	fmt.Print("Carrier Frequency: ")
	CARRIER_FREQ, _ = fmt.Scan()
	fmt.Print("Output Sample Rate: ")
	SAMPLE_RATE, _ = fmt.Scan()
}
*/


func main() {
	if checkGlobals() {
		data := []int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20}
		//data := []int{0, 0, 1, 1, 2, 2, 2, 3, 3, 4, 4, 4, 5, 5, 6, 6, 6, 7, 7, 8, 8, 8, 9, 9, 10, 10, 10, 11, 11, 12, 12, 12, 13, 13, 14, 14, 14, 15, 15, 16, 16, 16, 17, 17, 18, 18, 18, 19, 19, 19, 20, 20, 21, 21, 21, 22, 22, 23, 23, 23, 24, 24, 25, 25, 25, 26, 26, 27, 27, 27, 28, 28, 28, 29, 29, 30, 30, 30, 31, 31, 32, 32, 32, 33, 33, 34, 34, 34, 35, 35, 35, 36, 36, 37, 37, 37, 38, 38, 39, 39, 39, 40, 40, 40, 41, 41, 42, 42, 42, 43, 43, 43, 44, 44, 45, 45, 45, 46, 46, 46, 47, 47, 48, 48, 48, 49, 49, 49, 50, 50, 50, 51, 51, 52, 52, 52, 53, 53, 53, 54, 54, 54, 55, 55, 56, 56, 56, 57, 57, 57, 58, 58, 58, 59, 59, 59, 60, 60, 61, 61, 61, 62, 62, 62, 63, 63, 63, 64, 64, 64, 65, 65, 65, 66, 66, 66, 67, 67, 67, 68, 68, 68, 69, 69, 69, 70, 70, 70, 71, 71, 71, 72, 72, 72, 73, 73, 73, 74, 74, 74, 75, 75, 75, 76, 76, 76, 77, 77, 77, 78, 78, 78, 79, 79, 79, 79, 80, 80, 80, 81, 81, 81, 82, 82, 82, 83, 83, 83, 83, 84, 84, 84, 85, 85, 85, 86, 86, 86, 86, 87, 87, 87, 88, 88, 88, 88, 89, 89, 89, 90, 90, 90, 90, 91, 91, 91, 92, 92, 92, 92, 93, 93, 93, 93, 94, 94, 94, 95, 95, 95, 95, 96, 96, 96, 96, 97, 97, 97, 97, 98, 98, 98, 98, 99, 99, 99, 99, 100, 100, 100, 100, 101, 101, 101, 101, 102, 102, 102, 102, 103, 103, 103, 103, 103, 104, 104, 104, 104, 105, 105, 105, 105, 106, 106, 106, 106, 106, 107, 107, 107, 107, 107, 108, 108, 108, 108, 109, 109, 109, 109, 109, 110, 110, 110, 110, 110, 111, 111, 111, 111, 111, 112, 112, 112, 112, 112, 112, 113, 113, 113, 113, 113, 114, 114, 114, 114, 114, 114, 115, 115, 115, 115, 115, 115, 116, 116, 116, 116, 116, 116, 117, 117, 117, 117, 117, 117, 118, 118, 118, 118, 118, 118, 118, 119, 119, 119, 119, 119, 119, 119, 120, 120, 120, 120, 120, 120, 120, 120, 121, 121, 121, 121, 121, 121, 121, 121, 122, 122, 122, 122, 122, 122, 122, 122, 122, 123, 123, 123, 123, 123, 123, 123, 123, 123, 123, 124, 124, 124, 124, 124, 124, 124, 124, 124, 124, 124, 124, 125, 125, 125, 125, 125, 125, 125, 125, 125, 125, 125, 125, 125, 125, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 125, 125, 125, 125, 125, 125, 125, 125, 125, 125, 125, 125, 125, 125, 124, 124, 124, 124, 124, 124, 124, 124, 124, 124, 124, 124, 123, 123, 123, 123, 123, 123, 123, 123, 123, 123, 122, 122, 122, 122, 122, 122, 122, 122, 122, 121, 121, 121, 121, 121, 121, 121, 121, 120, 120, 120, 120, 120, 120, 120, 120, 119, 119, 119, 119, 119, 119, 119, 118, 118, 118, 118, 118, 118, 118, 117, 117, 117, 117, 117, 117, 117, 116, 116, 116, 116, 116, 116, 115, 115, 115, 115, 115, 115, 114, 114, 114, 114, 114, 113, 113, 113, 113, 113, 113, 112, 112, 112, 112, 112, 111, 111, 111, 111, 111, 110, 110, 110, 110, 110, 109, 109, 109, 109, 109, 108, 108, 108, 108, 108, 107, 107, 107, 107, 107, 106, 106, 106, 106, 105, 105, 105, 105, 105, 104, 104, 104, 104, 103, 103, 103, 103, 102, 102, 102, 102, 101, 101, 101, 101, 101, 100, 100, 100, 100, 99, 99, 99, 99, 98, 98, 98, 98, 97, 97, 97, 96, 96, 96, 96, 95, 95, 95, 95, 94, 94, 94, 94, 93, 93, 93, 92, 92, 92, 92, 91, 91, 91, 91, 90, 90, 90, 89, 89, 89, 89, 88, 88, 88, 87, 87, 87, 87, 86, 86, 86, 85, 85, 85, 84, 84, 84, 84, 83, 83, 83, 82, 82, 82, 81, 81, 81, 81, 80, 80, 80, 79, 79, 79, 78, 78, 78, 77, 77, 77, 76, 76, 76, 75, 75, 75, 75, 74, 74, 74, 73, 73, 73, 72, 72, 72, 71, 71, 71, 70, 70, 70, 69, 69, 69, 68, 68, 68, 67, 67, 67, 66, 66, 66, 65, 65, 65, 64, 64, 63, 63, 63, 62, 62, 62, 61, 61, 61, 60, 60, 60, 59, 59, 59, 58, 58, 58, 57, 57, 56, 56, 56, 55, 55, 55, 54, 54, 54, 53, 53, 52, 52, 52, 51, 51, 51, 50, 50, 50, 49, 49, 48, 48, 48, 47, 47, 47, 46, 46, 45, 45, 45, 44, 44, 44, 43, 43, 42, 42, 42, 41, 41, 41, 40, 40, 39, 39, 39, 38, 38, 38, 37, 37, 36, 36, 36, 35, 35, 34, 34, 34, 33, 33, 33, 32, 32, 31, 31, 31, 30, 30, 29, 29, 29, 28, 28, 28, 27, 27, 26, 26, 26, 25, 25, 24, 24, 24, 23, 23, 22, 22, 22, 21, 21, 20, 20, 20, 19, 19, 18, 18, 18, 17, 17, 17, 16, 16, 15, 15, 15, 14, 14, 13, 13, 13, 12, 12, 11, 11, 11, 10, 10, 9, 9, 9, 8, 8, 7, 7, 7, 6, 6, 5, 5, 5, 4, 4, 3, 3, 3, 2, 2, 1, 1, 1, 0, 0, -1, -1, -1, -2, -2, -3, -3, -3, -4, -4, -5, -5, -5, -6, -6, -7, -7, -7, -8, -8, -9, -9, -9, -10, -10, -11, -11, -11, -12, -12, -13, -13, -13, -14, -14, -15, -15, -15, -16, -16, -17, -17, -17, -18, -18, -18, -19, -19, -20, -20, -20, -21, -21, -22, -22, -22, -23, -23, -24, -24, -24, -25, -25, -26, -26, -26, -27, -27, -28, -28, -28, -29, -29, -29, -30, -30, -31, -31, -31, -32, -32, -33, -33, -33, -34, -34, -34, -35, -35, -36, -36, -36, -37, -37, -38, -38, -38, -39, -39, -39, -40, -40, -41, -41, -41, -42, -42, -42, -43, -43, -44, -44, -44, -45, -45, -45, -46, -46, -47, -47, -47, -48, -48, -48, -49, -49, -50, -50, -50, -51, -51, -51, -52, -52, -52, -53, -53, -54, -54, -54, -55, -55, -55, -56, -56, -56, -57, -57, -58, -58, -58, -59, -59, -59, -60, -60, -60, -61, -61, -61, -62, -62, -62, -63, -63, -63, -64, -64, -65, -65, -65, -66, -66, -66, -67, -67, -67, -68, -68, -68, -69, -69, -69, -70, -70, -70, -71, -71, -71, -72, -72, -72, -73, -73, -73, -74, -74, -74, -75, -75, -75, -75, -76, -76, -76, -77, -77, -77, -78, -78, -78, -79, -79, -79, -80, -80, -80, -81, -81, -81, -81, -82, -82, -82, -83, -83, -83, -84, -84, -84, -84, -85, -85, -85, -86, -86, -86, -87, -87, -87, -87, -88, -88, -88, -89, -89, -89, -89, -90, -90, -90, -91, -91, -91, -91, -92, -92, -92, -92, -93, -93, -93, -94, -94, -94, -94, -95, -95, -95, -95, -96, -96, -96, -96, -97, -97, -97, -98, -98, -98, -98, -99, -99, -99, -99, -100, -100, -100, -100, -101, -101, -101, -101, -101, -102, -102, -102, -102, -103, -103, -103, -103, -104, -104, -104, -104, -105, -105, -105, -105, -105, -106, -106, -106, -106, -107, -107, -107, -107, -107, -108, -108, -108, -108, -108, -109, -109, -109, -109, -109, -110, -110, -110, -110, -110, -111, -111, -111, -111, -111, -112, -112, -112, -112, -112, -113, -113, -113, -113, -113, -113, -114, -114, -114, -114, -114, -115, -115, -115, -115, -115, -115, -116, -116, -116, -116, -116, -116, -117, -117, -117, -117, -117, -117, -117, -118, -118, -118, -118, -118, -118, -118, -119, -119, -119, -119, -119, -119, -119, -120, -120, -120, -120, -120, -120, -120, -120, -121, -121, -121, -121, -121, -121, -121, -121, -122, -122, -122, -122, -122, -122, -122, -122, -122, -123, -123, -123, -123, -123, -123, -123, -123, -123, -123, -124, -124, -124, -124, -124, -124, -124, -124, -124, -124, -124, -124, -125, -125, -125, -125, -125, -125, -125, -125, -125, -125, -125, -125, -125, -125, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -125, -125, -125, -125, -125, -125, -125, -125, -125, -125, -125, -125, -125, -125, -124, -124, -124, -124, -124, -124, -124, -124, -124, -124, -124, -124, -123, -123, -123, -123, -123, -123, -123, -123, -123, -123, -122, -122, -122, -122, -122, -122, -122, -122, -122, -121, -121, -121, -121, -121, -121, -121, -121, -120, -120, -120, -120, -120, -120, -120, -120, -119, -119, -119, -119, -119, -119, -119, -118, -118, -118, -118, -118, -118, -118, -117, -117, -117, -117, -117, -117, -116, -116, -116, -116, -116, -116, -115, -115, -115, -115, -115, -115, -114, -114, -114, -114, -114, -114, -113, -113, -113, -113, -113, -112, -112, -112, -112, -112, -112, -111, -111, -111, -111, -111, -110, -110, -110, -110, -110, -109, -109, -109, -109, -109, -108, -108, -108, -108, -107, -107, -107, -107, -107, -106, -106, -106, -106, -106, -105, -105, -105, -105, -104, -104, -104, -104, -103, -103, -103, -103, -103, -102, -102, -102, -102, -101, -101, -101, -101, -100, -100, -100, -100, -99, -99, -99, -99, -98, -98, -98, -98, -97, -97, -97, -97, -96, -96, -96, -96, -95, -95, -95, -95, -94, -94, -94, -93, -93, -93, -93, -92, -92, -92, -92, -91, -91, -91, -90, -90, -90, -90, -89, -89, -89, -88, -88, -88, -88, -87, -87, -87, -86, -86, -86, -86, -85, -85, -85, -84, -84, -84, -83, -83, -83, -83, -82, -82, -82, -81, -81, -81, -80, -80, -80, -79, -79, -79, -79, -78, -78, -78, -77, -77, -77, -76, -76, -76, -75, -75, -75, -74, -74, -74, -73, -73, -73, -72, -72, -72, -71, -71, -71, -70, -70, -70, -69, -69, -69, -68, -68, -68, -67, -67, -67, -66, -66, -66, -65, -65, -65, -64, -64, -64, -63, -63, -63, -62, -62, -62, -61, -61, -61, -60, -60, -59, -59, -59, -58, -58, -58, -57, -57, -57, -56, -56, -56, -55, -55, -54, -54, -54, -53, -53, -53, -52, -52, -52, -51, -51, -50, -50, -50, -49, -49, -49, -48, -48, -48, -47, -47, -46, -46, -46, -45, -45, -45, -44, -44, -43, -43, -43, -42, -42, -42, -41, -41, -40, -40, -40, -39, -39, -39, -38, -38, -37, -37, -37, -36, -36, -35, -35, -35, -34, -34, -34, -33, -33, -32, -32, -32, -31, -31, -30, -30, -30, -29, -29, -28, -28, -28, -27, -27, -27, -26, -26, -25, -25, -25, -24, -24, -23, -23, -23, -22, -22, -21, -21, -21, -20, -20, -19, -19, -19, -18, -18, -18, -17, -17, -16, -16, -16, -15, -15, -14, -14, -14, -13, -13, -12, -12, -12, -11, -11, -10, -10, -10, -9, -9, -8, -8, -8, -7, -7, -6, -6, -6, -5, -5, -4, -4, -4, -3, -3, -2, -2, -2, -1, -1, 0, 0}
		bytes := toBytes(data)
		fmt.Println(bytes)
		
		bits := bytesToBinary(bytes)
		fmt.Println(bits)

		dataBitLength := SAMPLE_RATE / DATA_FREQ  //the length of 1 bit of data in samples, IF SAMPLE RATE IS 40000hz AND DATA RATE IS 100hz THEN IT SHOULD BE 400 SAMPLES WIDE
		carrierLen := len(bits) * dataBitLength
		carry := carrier(carrierLen)

		justData := modulate(carry, bits, dataBitLength)
		//fmt.Println(justData)

		finalData := addSilence(justData)
		//fmt.Println(finalData)

		file := wav.Create(finalData, 1, uint32(SAMPLE_RATE), 8)
		fmt.Println("Errors: ", file.Save("epicwin.wav"))
	}
}




































//BELOW IS THE CODE WHEN SILENCE WAS byte(127) INSTEAD OF byte(128)

/*
//DOES THE CARRIER_FREQ ACTUALLY NEED TO BE A MULTIPLE OF DATA_FREQ??? YOU SHOULD PROBABLY KNOW THAT
package main

import "wav"
import "fmt"


var DATA_FREQ = 100 //frequency that the data is sent in in hz
var CARRIER_FREQ = 4000 //frequency of carrier signal
var SAMPLE_RATE = 40000 //sample rate of the output, //MUST BE AN INTEGER MULTIPLE OF CARRIER_FREQ AND DATA_FREQ


func toBytes(data []int) []byte {
	bytes := make([]byte, len(data))
	for e, b := range data {
		bytes[e] = byte(b+127)
	}
	return bytes
}

func byte2Bits(b byte) []uint8 {
	mask := uint8(1)
	bits := make([]uint8, 8)
	for i := uint8(0); i < 8; i++ {
		if b & (mask << i) != 0 {
			bits[7 - i] = 1
		} //dont need to set any bits to 0 because make() initialized them to 0 by default
	}
	return bits
}

func bytesToBinary(bytes []byte) []byte {
	allBits := []byte{}
	currentBits := []byte{}
	for _, b := range bytes {
		currentBits = byte2Bits(b)
		currentBits = append([]byte{1}, currentBits...) //ADD START BIT
		currentBits = append(currentBits, 0) //ADD END BIT
		allBits = append(allBits, currentBits...)
	}
	return allBits
}

func carrier(length int) []byte {
	data := make([]byte, length)
	isHigh := false
	halfPeriod := SAMPLE_RATE / CARRIER_FREQ
	for i := 0; i < length; i++ {
		if i % halfPeriod == 0 { //every half a period
			isHigh = !isHigh //flip isHigh
		}
		if isHigh { //if its high, make the sample high
			data[i] = 255
		} else { //otherwise make it silence
			data[i] = 127
		}
	}
	return data
}

func modulate(carrier []byte, data []byte, bitLength int) []byte {
	carrierLoc := 0
	for _, b := range data {
		if b == 0 {
			for i := 0; i < bitLength; i++ { //make the next bit length of samples all 0
				carrier[carrierLoc + i] = 127
			}
		}
		carrierLoc += bitLength //move to the next data sample
	}
	return carrier
}

func addSilence(data []byte) []byte { //NEED TO ADD SILENCE BECAUSE WHEN I PLAY THE AUDIO IN WINDOWS MEDIA PLAYER IT FADES IN AND OUT
	silence := make([]byte, SAMPLE_RATE) //because sample rate is just the number of samples in a second
	for e := range silence {
		silence[e] = 127 //SILENCE IS 127 NOT 0!!!
	}
	data = append(silence, data...)
	return append(data, silence...)
}

func checkGlobals() bool {
	if SAMPLE_RATE % CARRIER_FREQ != 0 || SAMPLE_RATE % DATA_FREQ != 0 || CARRIER_FREQ % DATA_FREQ != 0 {
		fmt.Println("CARRIER OR DATA FREQUENCY IS NOT MULTIPLE OF SAMPLE_RATE")
		return false
	}
	return true
}


func main() {
	if checkGlobals() {
		data := []int{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20}
		//data := []int{0, 0, 1, 1, 2, 2, 2, 3, 3, 4, 4, 4, 5, 5, 6, 6, 6, 7, 7, 8, 8, 8, 9, 9, 10, 10, 10, 11, 11, 12, 12, 12, 13, 13, 14, 14, 14, 15, 15, 16, 16, 16, 17, 17, 18, 18, 18, 19, 19, 19, 20, 20, 21, 21, 21, 22, 22, 23, 23, 23, 24, 24, 25, 25, 25, 26, 26, 27, 27, 27, 28, 28, 28, 29, 29, 30, 30, 30, 31, 31, 32, 32, 32, 33, 33, 34, 34, 34, 35, 35, 35, 36, 36, 37, 37, 37, 38, 38, 39, 39, 39, 40, 40, 40, 41, 41, 42, 42, 42, 43, 43, 43, 44, 44, 45, 45, 45, 46, 46, 46, 47, 47, 48, 48, 48, 49, 49, 49, 50, 50, 50, 51, 51, 52, 52, 52, 53, 53, 53, 54, 54, 54, 55, 55, 56, 56, 56, 57, 57, 57, 58, 58, 58, 59, 59, 59, 60, 60, 61, 61, 61, 62, 62, 62, 63, 63, 63, 64, 64, 64, 65, 65, 65, 66, 66, 66, 67, 67, 67, 68, 68, 68, 69, 69, 69, 70, 70, 70, 71, 71, 71, 72, 72, 72, 73, 73, 73, 74, 74, 74, 75, 75, 75, 76, 76, 76, 77, 77, 77, 78, 78, 78, 79, 79, 79, 79, 80, 80, 80, 81, 81, 81, 82, 82, 82, 83, 83, 83, 83, 84, 84, 84, 85, 85, 85, 86, 86, 86, 86, 87, 87, 87, 88, 88, 88, 88, 89, 89, 89, 90, 90, 90, 90, 91, 91, 91, 92, 92, 92, 92, 93, 93, 93, 93, 94, 94, 94, 95, 95, 95, 95, 96, 96, 96, 96, 97, 97, 97, 97, 98, 98, 98, 98, 99, 99, 99, 99, 100, 100, 100, 100, 101, 101, 101, 101, 102, 102, 102, 102, 103, 103, 103, 103, 103, 104, 104, 104, 104, 105, 105, 105, 105, 106, 106, 106, 106, 106, 107, 107, 107, 107, 107, 108, 108, 108, 108, 109, 109, 109, 109, 109, 110, 110, 110, 110, 110, 111, 111, 111, 111, 111, 112, 112, 112, 112, 112, 112, 113, 113, 113, 113, 113, 114, 114, 114, 114, 114, 114, 115, 115, 115, 115, 115, 115, 116, 116, 116, 116, 116, 116, 117, 117, 117, 117, 117, 117, 118, 118, 118, 118, 118, 118, 118, 119, 119, 119, 119, 119, 119, 119, 120, 120, 120, 120, 120, 120, 120, 120, 121, 121, 121, 121, 121, 121, 121, 121, 122, 122, 122, 122, 122, 122, 122, 122, 122, 123, 123, 123, 123, 123, 123, 123, 123, 123, 123, 124, 124, 124, 124, 124, 124, 124, 124, 124, 124, 124, 124, 125, 125, 125, 125, 125, 125, 125, 125, 125, 125, 125, 125, 125, 125, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 127, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 126, 125, 125, 125, 125, 125, 125, 125, 125, 125, 125, 125, 125, 125, 125, 124, 124, 124, 124, 124, 124, 124, 124, 124, 124, 124, 124, 123, 123, 123, 123, 123, 123, 123, 123, 123, 123, 122, 122, 122, 122, 122, 122, 122, 122, 122, 121, 121, 121, 121, 121, 121, 121, 121, 120, 120, 120, 120, 120, 120, 120, 120, 119, 119, 119, 119, 119, 119, 119, 118, 118, 118, 118, 118, 118, 118, 117, 117, 117, 117, 117, 117, 117, 116, 116, 116, 116, 116, 116, 115, 115, 115, 115, 115, 115, 114, 114, 114, 114, 114, 113, 113, 113, 113, 113, 113, 112, 112, 112, 112, 112, 111, 111, 111, 111, 111, 110, 110, 110, 110, 110, 109, 109, 109, 109, 109, 108, 108, 108, 108, 108, 107, 107, 107, 107, 107, 106, 106, 106, 106, 105, 105, 105, 105, 105, 104, 104, 104, 104, 103, 103, 103, 103, 102, 102, 102, 102, 101, 101, 101, 101, 101, 100, 100, 100, 100, 99, 99, 99, 99, 98, 98, 98, 98, 97, 97, 97, 96, 96, 96, 96, 95, 95, 95, 95, 94, 94, 94, 94, 93, 93, 93, 92, 92, 92, 92, 91, 91, 91, 91, 90, 90, 90, 89, 89, 89, 89, 88, 88, 88, 87, 87, 87, 87, 86, 86, 86, 85, 85, 85, 84, 84, 84, 84, 83, 83, 83, 82, 82, 82, 81, 81, 81, 81, 80, 80, 80, 79, 79, 79, 78, 78, 78, 77, 77, 77, 76, 76, 76, 75, 75, 75, 75, 74, 74, 74, 73, 73, 73, 72, 72, 72, 71, 71, 71, 70, 70, 70, 69, 69, 69, 68, 68, 68, 67, 67, 67, 66, 66, 66, 65, 65, 65, 64, 64, 63, 63, 63, 62, 62, 62, 61, 61, 61, 60, 60, 60, 59, 59, 59, 58, 58, 58, 57, 57, 56, 56, 56, 55, 55, 55, 54, 54, 54, 53, 53, 52, 52, 52, 51, 51, 51, 50, 50, 50, 49, 49, 48, 48, 48, 47, 47, 47, 46, 46, 45, 45, 45, 44, 44, 44, 43, 43, 42, 42, 42, 41, 41, 41, 40, 40, 39, 39, 39, 38, 38, 38, 37, 37, 36, 36, 36, 35, 35, 34, 34, 34, 33, 33, 33, 32, 32, 31, 31, 31, 30, 30, 29, 29, 29, 28, 28, 28, 27, 27, 26, 26, 26, 25, 25, 24, 24, 24, 23, 23, 22, 22, 22, 21, 21, 20, 20, 20, 19, 19, 18, 18, 18, 17, 17, 17, 16, 16, 15, 15, 15, 14, 14, 13, 13, 13, 12, 12, 11, 11, 11, 10, 10, 9, 9, 9, 8, 8, 7, 7, 7, 6, 6, 5, 5, 5, 4, 4, 3, 3, 3, 2, 2, 1, 1, 1, 0, 0, -1, -1, -1, -2, -2, -3, -3, -3, -4, -4, -5, -5, -5, -6, -6, -7, -7, -7, -8, -8, -9, -9, -9, -10, -10, -11, -11, -11, -12, -12, -13, -13, -13, -14, -14, -15, -15, -15, -16, -16, -17, -17, -17, -18, -18, -18, -19, -19, -20, -20, -20, -21, -21, -22, -22, -22, -23, -23, -24, -24, -24, -25, -25, -26, -26, -26, -27, -27, -28, -28, -28, -29, -29, -29, -30, -30, -31, -31, -31, -32, -32, -33, -33, -33, -34, -34, -34, -35, -35, -36, -36, -36, -37, -37, -38, -38, -38, -39, -39, -39, -40, -40, -41, -41, -41, -42, -42, -42, -43, -43, -44, -44, -44, -45, -45, -45, -46, -46, -47, -47, -47, -48, -48, -48, -49, -49, -50, -50, -50, -51, -51, -51, -52, -52, -52, -53, -53, -54, -54, -54, -55, -55, -55, -56, -56, -56, -57, -57, -58, -58, -58, -59, -59, -59, -60, -60, -60, -61, -61, -61, -62, -62, -62, -63, -63, -63, -64, -64, -65, -65, -65, -66, -66, -66, -67, -67, -67, -68, -68, -68, -69, -69, -69, -70, -70, -70, -71, -71, -71, -72, -72, -72, -73, -73, -73, -74, -74, -74, -75, -75, -75, -75, -76, -76, -76, -77, -77, -77, -78, -78, -78, -79, -79, -79, -80, -80, -80, -81, -81, -81, -81, -82, -82, -82, -83, -83, -83, -84, -84, -84, -84, -85, -85, -85, -86, -86, -86, -87, -87, -87, -87, -88, -88, -88, -89, -89, -89, -89, -90, -90, -90, -91, -91, -91, -91, -92, -92, -92, -92, -93, -93, -93, -94, -94, -94, -94, -95, -95, -95, -95, -96, -96, -96, -96, -97, -97, -97, -98, -98, -98, -98, -99, -99, -99, -99, -100, -100, -100, -100, -101, -101, -101, -101, -101, -102, -102, -102, -102, -103, -103, -103, -103, -104, -104, -104, -104, -105, -105, -105, -105, -105, -106, -106, -106, -106, -107, -107, -107, -107, -107, -108, -108, -108, -108, -108, -109, -109, -109, -109, -109, -110, -110, -110, -110, -110, -111, -111, -111, -111, -111, -112, -112, -112, -112, -112, -113, -113, -113, -113, -113, -113, -114, -114, -114, -114, -114, -115, -115, -115, -115, -115, -115, -116, -116, -116, -116, -116, -116, -117, -117, -117, -117, -117, -117, -117, -118, -118, -118, -118, -118, -118, -118, -119, -119, -119, -119, -119, -119, -119, -120, -120, -120, -120, -120, -120, -120, -120, -121, -121, -121, -121, -121, -121, -121, -121, -122, -122, -122, -122, -122, -122, -122, -122, -122, -123, -123, -123, -123, -123, -123, -123, -123, -123, -123, -124, -124, -124, -124, -124, -124, -124, -124, -124, -124, -124, -124, -125, -125, -125, -125, -125, -125, -125, -125, -125, -125, -125, -125, -125, -125, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -127, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -126, -125, -125, -125, -125, -125, -125, -125, -125, -125, -125, -125, -125, -125, -125, -124, -124, -124, -124, -124, -124, -124, -124, -124, -124, -124, -124, -123, -123, -123, -123, -123, -123, -123, -123, -123, -123, -122, -122, -122, -122, -122, -122, -122, -122, -122, -121, -121, -121, -121, -121, -121, -121, -121, -120, -120, -120, -120, -120, -120, -120, -120, -119, -119, -119, -119, -119, -119, -119, -118, -118, -118, -118, -118, -118, -118, -117, -117, -117, -117, -117, -117, -116, -116, -116, -116, -116, -116, -115, -115, -115, -115, -115, -115, -114, -114, -114, -114, -114, -114, -113, -113, -113, -113, -113, -112, -112, -112, -112, -112, -112, -111, -111, -111, -111, -111, -110, -110, -110, -110, -110, -109, -109, -109, -109, -109, -108, -108, -108, -108, -107, -107, -107, -107, -107, -106, -106, -106, -106, -106, -105, -105, -105, -105, -104, -104, -104, -104, -103, -103, -103, -103, -103, -102, -102, -102, -102, -101, -101, -101, -101, -100, -100, -100, -100, -99, -99, -99, -99, -98, -98, -98, -98, -97, -97, -97, -97, -96, -96, -96, -96, -95, -95, -95, -95, -94, -94, -94, -93, -93, -93, -93, -92, -92, -92, -92, -91, -91, -91, -90, -90, -90, -90, -89, -89, -89, -88, -88, -88, -88, -87, -87, -87, -86, -86, -86, -86, -85, -85, -85, -84, -84, -84, -83, -83, -83, -83, -82, -82, -82, -81, -81, -81, -80, -80, -80, -79, -79, -79, -79, -78, -78, -78, -77, -77, -77, -76, -76, -76, -75, -75, -75, -74, -74, -74, -73, -73, -73, -72, -72, -72, -71, -71, -71, -70, -70, -70, -69, -69, -69, -68, -68, -68, -67, -67, -67, -66, -66, -66, -65, -65, -65, -64, -64, -64, -63, -63, -63, -62, -62, -62, -61, -61, -61, -60, -60, -59, -59, -59, -58, -58, -58, -57, -57, -57, -56, -56, -56, -55, -55, -54, -54, -54, -53, -53, -53, -52, -52, -52, -51, -51, -50, -50, -50, -49, -49, -49, -48, -48, -48, -47, -47, -46, -46, -46, -45, -45, -45, -44, -44, -43, -43, -43, -42, -42, -42, -41, -41, -40, -40, -40, -39, -39, -39, -38, -38, -37, -37, -37, -36, -36, -35, -35, -35, -34, -34, -34, -33, -33, -32, -32, -32, -31, -31, -30, -30, -30, -29, -29, -28, -28, -28, -27, -27, -27, -26, -26, -25, -25, -25, -24, -24, -23, -23, -23, -22, -22, -21, -21, -21, -20, -20, -19, -19, -19, -18, -18, -18, -17, -17, -16, -16, -16, -15, -15, -14, -14, -14, -13, -13, -12, -12, -12, -11, -11, -10, -10, -10, -9, -9, -8, -8, -8, -7, -7, -6, -6, -6, -5, -5, -4, -4, -4, -3, -3, -2, -2, -2, -1, -1, 0, 0}
		bytes := toBytes(data)
		fmt.Println(bytes)
		
		bits := bytesToBinary(bytes)
		fmt.Println(bits)

		dataBitLength := SAMPLE_RATE / DATA_FREQ  //the length of 1 bit of data in samples, IF SAMPLE RATE IS 40000hz AND DATA RATE IS 100hz THEN IT SHOULD BE 400 SAMPLES WIDE
		carrierLen := len(bits) * dataBitLength
		carry := carrier(carrierLen)

		justData := modulate(carry, bits, dataBitLength)
		//fmt.Println(justData)

		finalData := addSilence(justData)
		//fmt.Println(finalData)

		file := wav.Create(finalData, 1, uint32(SAMPLE_RATE), 8)
		fmt.Println("Errors: ", file.Save("epicwin.wav"))
	}
}
*/