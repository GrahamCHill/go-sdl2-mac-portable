# Go SDL2 Portable macOS Binary
A simple application that shows how to build an SDL2 mac app with the Go Language.  
This assumes you have followed all setup available from https://github.com/veandco/go-sdl2

## License
Please be aware while this project is licensed under the MIT license, the projects for SDL2 and go-sdl2 have other 
licenses, while I waive the right to keep the copy of this project's MIT license with the shareable project, you still 
have to include and abide by the other projects requirements when sharing.


## Steps to build on macOS
### Building SDL2 for universal builds
1. Build SDL2 from source
   1. a.	SDL2 is available from “ https://github.com/libsdl-org/SDL * “ you will need cmake and make to build it, 
   open the terminal in the downloaded SDL2 folder, and run the following to build the universal library.
   
```bash
mkdir build && cd build
cmake .. -DCMAKE_OSX_ARCHITECTURES="arm64;x86_64"
make
otool -L libSDL2-2.0.0.dylib
```  
The above may take a while but when it is done you should get an output telling you what native libraries 
libSDL2-2.0.0.dylib is linking against (this file should also be fully portable.)	

\* make sure you are using SDL2 and not SDL3, you will have to swap to the correct branch to make sure you are using the right one

2. Copy the newly built dylib (libSDL2-2.0.0.dylib) to usr/local/lib as the Go build system will need this file (keep a 
copy available as you will need it later.)

### Building your SDL App for both x86-64 and arm64
Now once you have written your code you can do the following to build and link your files, this won't go into building 
an app bundle as is the standard way to share mac binaries.  

This guide uses the names like awesomeOutputARM64, and awesomeOutputAMD64, you can change these names to whatever works 
for you or rename your files at the end so you can just copy-paste the commands (I may create a script that automates 
these steps at a later date.)

#### AMD64 (x86-64) binary
For AMD64 in your Go project folder run the following command
```bash
GOARCH=amd64 GOOS=darwin CGO_ENABLED=1 go build -o awesomeOutputAMD64 .
```

This will output an AMD64 binary from your Go code

#### ARM64 binary
For regular ARM64 binaries, you can either rely on the regular go build command, or use the below to keep the process 
simple

```bash
GOARCH=arm64 GOOS=darwin CGO_ENABLED=1 go build -o awesomeOutputARM64 .
```

Now with the two separate binaries, if the builds succeeded, you can run the following commands to confirm that they are 
built for AMD64 and ARM64 respectively.
```bash
lipo -info awesomeOutputAMD64
lipo -info awesomeOutputARM64
```

you should see an output of
```text
Non-fat file: awesomeOutputAMD64 is architecture: x86_64
Non-fat file: awesomeOutputARM64 is architecture: arm64
```

If all has worked you can move onto the next steps

### Building the Universal Go Binary
At this point with the two executables you can run the following command to create a fat/universal binary 
```bash
lipo -create -output ./awesomeOutputUniversal ./awesomeOutputARM64 ./awesomeOutputAMD64 
lipo -info awesomeOutputUniversal
```

you should get the following output, as well as a new binary alongside the two separate architecture binaries
```text
Architectures in the fat file: awesomeOutputUniversal are: x86_64 arm64
```

### Making it Portable

Now to make it portable (this is not an app bundle, but the same logic can be used, just alter commands as needed, where
needed.)
First run the following command in your terminal
```bash
otool -L awesomeOutputUniversal
```

This should output something like the following (if you are using other SDL2 libraries, or just other libraries in 
general, you can modify the commands I have provided to create a portable library, and the output below will differ.)
```text
awesomeOutputUniversal (architecture x86_64)
    @rpath/libSDL2-2.0.0.dylib (compatibility version 3101.0.0, current version 3101.0.0)
    /usr/lib/libSystem.B.dylib (compatibility version 1.0.0, current version 1351.0.0)
    /usr/lib/libresolv.9.dylib (compatibility version 1.0.0, current version 1.0.0)
awesomeOutputUniversal (architecture arm64)
    @rpath/libSDL2-2.0.0.dylib (compatibility version 3101.0.0, current version 3101.0.0)
    /System/Library/Frameworks/CoreFoundation.framework/Versions/A/CoreFoundation (compatibility version 150.0.0, current version 3107.0.0)
    /usr/lib/libresolv.9.dylib (compatibility version 1.0.0, current version 1.0.0)
    /usr/lib/libSystem.B.dylib (compatibility version 1.0.0, current version 1351.0.0)
```
To make your application portable you are only interested in the lines that say things like @rpath, or have your 
username in them. For SDL2 in this application the line is `@rpath/libSDL2-2.0.0.dylib`  

You can then run the following command to finally make your universal application portable
```bash
install_name_tool -change "@rpath/libSDL2-2.0.0.dylib" "@executable_path/./libSDL2-2.0.0.dylib" awesomeOutputUniversal
```

Now you can share your Go application that uses SDL2 with other macOS users

### Notes on app bundles 
To make an app bundle portable, most of the above does apply with minor alterations, the most important of which is the 
last command, which will change to
```bash
install_name_tool -change "@rpath/libSDL2-2.0.0.dylib" "@executable_path/../Frameworks/libSDL2-2.0.0.dylib" awesomeOutputUniversal
```

Since your typical macOS app bundle will have the following structure
```text
-awesomeApp.app
    -info.plist
    -Frameworks
        -libSDL2-2.0.0.dylib
    -MacOS
        -awesomeAppUniversal
    -Resources
        -awesomeAppIcon.icns
        -other app resources like images
```
Typically a mac app bundle would have the same name for the file in the MacOS folder and the app bundle, but just to 
avoid confusion with the commands I have named them differently. 