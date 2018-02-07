#DON'T USE __iter__ AND next, INSTEAD USE __getitem__ SO THAT YOU DON'T GET STUCK HALFWAY THROUGH ITERATING AND STUFF
    #__getitem__ is automatically called when doing for in loops, FOR LOOPS ONLY KNOW WHEN TO STOP BY CHECKING IF __getitem__ RAISES StopIteration()
    
#MIGHT WANNA NOT READ ALL THE FRAMES RIGHT AWAY, INSTEAD ALWAYS CALL wave.readframes OR SOMETHING

#~~~SHOULD REALLY FIGURE OUT WHAT EXACTLY 0 VOLUME IS (its probably 127 but I could have really screwed that up)
    #I THINK SILENCE SHOULD ACTUALLY BE 128 HOLY MOLY (since the lowest amplitude should be -128 which should map to the null byte)
#NEED TO FIGURE OUT IF w.readframes()[-1] IS DIFFERENT THAN open(file_loc).read()[44:][-1]
    #MAKE SURE TO TEST 16 BIT 2 CHANNELS STUFF

#!!!SILENCE EXPORTED BY ABLETON HAS A BUNCH OF RANDOM JUNK IN IT FOR SOME REASON

import wave


class audio:
    '''
    A class that represents an audio clip
    file_loc = path to a .wav file to read audio from
    
    Once open, you can access elements of the wave file like so:
        a = audio(file_loc)
        a[0]  #returns first sample
        a[-1]  #returns last sample
        for s in a:  #iterates through samples
            print s
            
    Samples are returned in a Tuple, with the amplitude of the first channel
    in the 0th element of the Tuple

    Other methods include:
        a.sample_rate()  #returns sample rate of the audio
        a.channels()  #returns the number of channels the audio has
        a.byte_depth()  #returns the bit depth of the audio in bytes
    '''
    
    def __init__(self, file_loc):
        self.wav = wave.open(file_loc)
        self._raw = self.wav.readframes(len(self))
        self.wav.close()
        
        self.sample_rate = self.wav.getframerate
        self.channels = self.wav.getnchannels
        self.byte_depth = self.wav.getsampwidth  #gets bit depth in bytes
        self._check_raw()  #make sure the wave file is formatted properly

    def __getitem__(self, index):
        if index < 0:  #so that index of -1 accesses the last item
            index = len(self) + index
        if index < len(self) and index >= 0:
            totals = []
            for c in range(self.channels()):
                chan_total = 0
                for b in range(self.byte_depth()):
                    chan_total <<= 8  #shift for next byte
                    byte = self._get_byte(index, c, b)  #access the proper byte
                    chan_total += ord(byte)  #add the value of the byte
                chan_total -= 2**(self.byte_depth()*8-1)  #convert to signed (by subtracting 128 if 8 bit)
                totals.append(chan_total)
            return tuple(totals)
        else:
            raise StopIteration()

    def __len__(self):
        return self.wav.getnframes()

    def _check_raw(self):  #checks if wave file is properly formatted
        if len(self._raw)%(self.byte_depth()*self.channels()) != 0:
            raise Exception("Wave file is not formatted correctly!")

    def _get_byte(self, sample_index, chan_num, byte_num):  #returns the specified byte of the specified channel of the specified sample
        return self._raw[sample_index * self.channels() * self.byte_depth()  #the index of the 0th byte of the sample
                         + chan_num * self.byte_depth()  #move to the 0th byte of the correct channel of the sample
                         + byte_num]  #move to the correct byte in the channel




if __name__ == "__main__":
    le = audio('noidbeamDataTrue.wav')
    with open('noidbeamDataTrue.wav', 'rb') as f:
        ze = f.read()[44:]

    print len(ze)
    print len(le)

    x = le[10001]
    y = ze[10001*4:10001*4+4]
