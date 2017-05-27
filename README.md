# File Arranger
To group files by file modification time.
If you are a photographer like me, you might have a trouble to manage thousands and thousands of photos, one way is to group those photos by dates, my prefered file structure is like:
```
year
   month
      day
         photo1
         photo2
```
But manually move files are boring, that's why I wrote this.
## Build executable
```
go build
```

## Usage
1. plug your SD card, let's say it's under `/dev/xcard`
2. run the following command:
```
./file-arranger -from /dev/xcard/camera -to /photos
```

