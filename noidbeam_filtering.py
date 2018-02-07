import scipy.signal as signal
import wave

#!!!DON'T FORGET TO CHANGE SAMPLING_FREQ IF YOU NEED TO

SAMPLING_FREQ = 44100 #CHANGE THIS IF YOU'RE READING DIFFERENT SAMP RATE YOU IMBECIL

def bandpass(data, low, high, order = 5):
    nyquist = SAMPLING_FREQ * 0.5
    low = low/nyquist
    high = high/nyquist
    b, a = signal.butter(order, [low, high], btype = 'band') #no idea what b and a are, but this makes a bandpass filter
    #return lfilter(b, a, data) #PHASE SHIFTS THE OUTPUT but you can use it in real time apparently
    ary = signal.filtfilt(b, a, data) #apply the butter bandpass filter to data, WITHOUT PHASE SHIFTS, to achieve this it filters forward and backwards
    return list(ary)



if __name__ == '__main__':

    def num2chr(n):
        return chr(int(n) + 127)
    
    def getData(path): #returns a list of all the data in noidbeamData.wav (data list is a list of amplitudes from 1 to 127)
        f = open(path, 'rb')
        raw = f.read()
        f.close()
        raw = raw[44:] #cut first 44 bytes of the wav file
        return [ord(a) for a in raw]

    def writeData(path, data):
        f = wave.open(path, 'w')
        f.setnchannels(1)
        f.setsampwidth(1)
        f.setframerate(SAMPLING_FREQ)
        f.writeframes(''.join([num2chr(d) for d in data]))
        f.close()

        
    data = getData('noidbeamDataTrue_8bit_1channel_STILL44100HZ.wav')
    print data[:100]
    data = bandpass(data, 1500, 2500)
    print data[:100]
    print max(data)
    writeData('noidbeamDataTrue_8bit_1channel_STILL44100HZ_filtered.wav', data)

