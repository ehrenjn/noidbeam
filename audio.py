#MIGHT WANNA NOT READ ALL THE FRAMES RIGHT AWAY, INSTEAD ALWAYS CALL wave.readframes OR SOMETHING
    #even a 5 minute wav file only takes like a second to read so it doesn't really matter

#!!!SILENCE EXPORTED BY ABLETON HAS A BUNCH OF RANDOM JUNK IN IT FOR SOME REASON
    #probably dithering

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
        self._channels = self.wav.getnchannels()
        self.byte_depth = self.wav.getsampwidth  #gets bit depth in bytes
        self._byte_depth = self.wav.getsampwidth()
        self._check_raw()  #make sure the wave file is formatted properly

    def __getitem__(self, index):
        if index < 0:  #so that index of -1 accesses the last item
            index = len(self) + index
        if index < len(self) and index >= 0:
            totals = []
            for c in range(self._channels):
                chan_total = 0
                for b in range(self._byte_depth):
                    chan_total <<= 8  #shift for next byte
                    byte = self._get_byte(index, c, b)  #access the proper byte
                    chan_total += ord(byte)  #add the value of the byte
                chan_total = self._convert_to_signed(chan_total)  #convert to signed int
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
        return self._raw[sample_index * self._channels * self._byte_depth  #the index of the 0th byte of the sample
                         + chan_num * self._byte_depth   #move to the 0th byte of the correct channel of the sample
                         + (self._byte_depth - 1)  #move to the last byte of the correct channel of the sample
                         - byte_num]  #move back to the correct byte in the channel, BYTES ARE STORED IN LITTLE ENDIAN ORDER FOR SOME UNGODLY REASON

    def _convert_to_signed(self, n):  #converts n to a signed number
        if self.byte_depth() == 1:  #if 8 bit then the null byte is just -128 and \xff is 127
            n -= 2**(self._byte_depth*8-1)  #subtract 128 if 8 bit
        else:  #IF IT ISN'T 8 BIT THEN THE AUDIO DATA IS STORED IN 2'S COMPLEMENT ENCODING
            shift_dist = (self._byte_depth*8) - 1
            sign = n >> shift_dist
            n -= (sign << shift_dist) * 2  #I could have done n -= n >> shift_dist << (shift_dist + 1)  but this is way more readable
        return n
