# PI Video Builder

This program is used to generate a video sequence of the digits of PI based on a series of input video and input text list

## Requirements

1. Windows 10 OS
1. 32GB Ram (will work with less RAM, but will be much slower)
1. Make sure you have FFmpeg installed and accessible from command line
1. Folder containing only `.MP4` video clips for digits 0-9. Clips will be read in alphanumeric order (IE: First clip 1. alphanumerically will be assumed to be the digit 0)
1. A `.TXT` file containing digits of PI. PVB will automatically remove all non-numeric data

## How to run

1. Place `pi-video-builder.exe` into the root project directory
1. With the command line pointed to the project directory, run `.\pi-video-exitor.exe [src_folder\path] [path\pi_text_file.txt] [path\to\output_folder]` (without [ ] braces)

## Pausing, Resuming and Restarting

If you want to pause the render, you can do so easily by: 

1. Exiting out of the command prompt window
1. Deleting the most recent video output. If this is not deleted there will be an error on next render as it will not overwrite existing video.
1. To resume, simply run the same command you ran to originally run the program. PVB will re-start where it left off
1. **WARNING:** Do not delete `log.txt` file. If this is deleted PVB will be unable to resume where it left off

If you would like to restart the process from scratch: 

1. Delete all generated videos
1. Delete the `log.txt` file from the program directory
1. Run the original command