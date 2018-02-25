# noidbeam
File transfer using sound (work in progress)  

noidbeam.go converts a file into a series of audio pulses and saves the noise to a .wav file. To send the file to another computer the user plays the .wav out loud on the first computer (as loud as possible without clipping) and records the audio with the 2nd computer. Once the audio is recorded and saved as a .wav file, noidbeam_recieve.py should be able to read it and print out the contents of the file.  
The next thing I want to add is frequency multiplexing to increase the efficiency of the file sending.
