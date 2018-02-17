//DOES THE CARRIER_FREQ ACTUALLY NEED TO BE A MULTIPLE OF DATA_FREQ??? YOU SHOULD PROBABLY KNOW THAT
package main

import "wav"
import "fmt"
import "os"


var DATA_FREQ = 100 //frequency that the data is sent in in hz
var CARRIER_FREQ = 4000 //frequency of carrier signal
var SAMPLE_RATE = 40000 //sample rate of the output, //MUST BE AN INTEGER MULTIPLE OF CARRIER_FREQ AND DATA_FREQ
var DATA_BIT_LENGTH = SAMPLE_RATE / DATA_FREQ //the length of 1 bit of data in samples, IF SAMPLE RATE IS 40000hz AND DATA RATE IS 100hz THEN IT SHOULD BE 400 SAMPLES WIDE

var SILENCE = byte(128) //NOT 0 OR 127; its 128 because amplitude is -128 to 127, NOT -127 to 127


func readFile(fileLocation string) ([]byte) { //reads file at fileLocation to a byte slice
	f, err := os.Open(fileLocation)
	if err != nil {
		fmt.Println(fileLocation, "CAN'T BE OPENED!")
		os.Exit(69)
	}
	info, _ := f.Stat()
	size := info.Size()
	buffer := make([]byte, size)
	f.Read(buffer)
	f.Close()
	return buffer
}

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

func modulate(carrier []byte, data []byte) []byte {
	carrierLoc := 0
	for _, b := range data {
		if b == 0 {
			for i := 0; i < DATA_BIT_LENGTH; i++ { //make the next bit length of samples all silent
				carrier[carrierLoc + i] = SILENCE
			}
		}
		carrierLoc += DATA_BIT_LENGTH //move to the next data sample
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

func addCalibrationBits(data []byte) []byte { //add bits at the beginning to calibrate the noise threshold on the recieving end
	calibration := []byte{1,0,1,0,1,0,1,0,1,0}
	carry := carrier(DATA_BIT_LENGTH * len(calibration))
	data = append(modulate(carry, calibration), data...)
	return data
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
		bytes := readFile("noidbeamData.txt")
		fmt.Println(bytes)
		fmt.Println("Read", len(bytes), "bytes")
		
		bits := bytesToBinary(bytes)
		fmt.Println(bits)
 
		carrierLen := len(bits) * DATA_BIT_LENGTH
		carry := carrier(carrierLen)

		justData := modulate(carry, bits)
		//fmt.Println(justData)

		dataWithCalibration := addCalibrationBits(justData)
		//fmt.Println(dataWithCalibration)

		finalData := addSilence(dataWithCalibration)
		//fmt.Println(finalData)


		file := wav.Create(finalData, 1, uint32(SAMPLE_RATE), 8)
		fmt.Println("Errors: ", file.Save("epicwin.wav"))
	}
}
