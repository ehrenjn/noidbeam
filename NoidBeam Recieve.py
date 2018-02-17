
#import noidbeam_filtering as filtering
from audio import audio


DATA = audio('noibeamTrueDataCalibrated.wav')
DATA_BIT_RATE = 100
DATA_SAMPLE_WIDTH = DATA.sample_rate() / DATA_BIT_RATE



def noise_floor(data): #finds the average amplitude when data is not being sent
    quarter_second = data.sample_rate() / 4 #number of samples in 1 quarter of second
    return sum(abs(data[d][0]) for d in xrange(quarter_second, quarter_second*2)) / quarter_second #returns the average of the 2nd quarter second of recording (FIRST SECOND OF RECORDING SHOULD BE SILENT)

#AVERAGE_THRESHOLD =  18
AVERAGE_THRESHOLD = 0 #if the average amplitude of a bunch of samples is above this threshold, those samples are considered a 1
MAX_AMPLITUDE = 0 #theoretical max amplitude of a 1, if a sample is 50% of this number or more it's probably in a 1
def determine_bit(samples): #turns a list of samples into a bit (figures out if the sample list represents a 1 or a 0)
    return int((reduce(lambda x,y: abs(x) + abs(y), samples)/len(samples)) > AVERAGE_THRESHOLD)
    
def next_bits(data, start, num_bits): #returns the specified number of bits from data starting at start
    end = start + DATA_SAMPLE_WIDTH * num_bits
    byte = [data[d][0] for d in xrange(start, end)]
    bits = [determine_bit(byte[DATA_SAMPLE_WIDTH*b : DATA_SAMPLE_WIDTH*(b+1)]) for b in range(num_bits)]
    return bits

def calibrate_threshold(data, start):
    global AVERAGE_THRESHOLD
    print "Calibrating threshold..."
    bits = []
    threshold_change = 1 if data.byte_depth() == 1 else 2*(10**data.byte_depth())  #amt the threshold will be changed every attempt to read calibration bits
    while bits != [1,0] * 5:
        AVERAGE_THRESHOLD += threshold_change
        bits = next_bits(data, start, 10)
        print bits
        if bits == [0]*10: #AVERAGE_THRESHOLD >= 2**(data.byte_depth()*8):
            print "CALIBRATION BITS COULDN'T BE READ!"
            exit()
    print "Average threshold: " + str(AVERAGE_THRESHOLD)
    return bits

def find_start_bit(data):
    global MAX_AMPLITUDE
    print "Noise floor: " + str(noise_floor(data))
    print "Finding start bit..."
    beginning_len = data.sample_rate() * 5
    if beginning_len > len(data): #if the data is less than 5 seconds long just go through all the data to find the start
        beginning_len = len(data)
    first_5_seconds = [abs(data[d][0]) for d in xrange(beginning_len)]
    max_amp = max(first_5_seconds)
    MAX_AMPLITUDE = max_amp
    for i, d in enumerate(first_5_seconds):
        if d >= max_amp/3.0: #get the first sample thats at least 1/3rd of the max of the first 5 seconds
            print "Start bit @ sample " + str(i)
            return i
    print "COULDN'T FIND START BIT"
    exit()

def realign(start_of_bit, data):
    start_of_search = start_of_bit - (DATA_SAMPLE_WIDTH / 2) #go half a bit back and start looking for the start of this bit
    for b in xrange(start_of_search, start_of_search + DATA_SAMPLE_WIDTH):
        if abs(data[b][0]) >= MAX_AMPLITUDE * 0.5:
            amt_moved = start_of_bit - b
            if amt_moved >= 0:
                print "moved " + str(amt_moved) + " samples backward"
            else:
                print "moved " + str(abs(amt_moved)) + " samples forward"
            return b
    print "COULDN'T FIND START BIT"
    return start_of_bit

def bits_to_bytes(bits): #converts a list of bits to a list of bytes
    if len(bits)%8 != 0:
        print "WARNING: RECIEVED DATA IS INCOMPLETE (num of bits is not a multiple of 8)" #dont think its actually ever gonna print this cause it only checks for start and stop every 8 bits
    byte_list = []
    current_byte = 0
    for i, b in enumerate(bits):
        current_byte <<= 1
        current_byte += b
        if i % 8 == 7: #if a full bit has been found
            byte_list.append(current_byte)
            current_byte = 0
    return byte_list

def bytes_to_str(byte_list): #converts a list of bytes to a string
    return ''.join(chr(b) for b in byte_list)
    


#data = [int(d) for d in filtering.bandpass(data, 1500, 2500)]
data_bits = []
current_sample = find_start_bit(DATA)
calibrate_threshold(DATA, current_sample)
current_sample += DATA_SAMPLE_WIDTH * 10 #move ahead past all the calibration bits
print "\nReading data..."
while next_bits(DATA, current_sample, 1)[0]: #while the next start bit is 1
    current_sample += DATA_SAMPLE_WIDTH * 1 #move ahead a bit
    data_bits += next_bits(DATA, current_sample, 8) #nab the data byte
    current_sample += DATA_SAMPLE_WIDTH * 9 #move ahead a byte PLUS the stop bit
    if len(data_bits) % 80 == 0:
        print "\t" + str(len(data_bits)/8) + " bytes read"
    current_sample = realign(current_sample, DATA)
    
print data_bits
byte_list = bits_to_bytes(data_bits)
print byte_list
print bytes_to_str(byte_list)
