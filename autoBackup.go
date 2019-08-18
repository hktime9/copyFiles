package main

import(
	"fmt"
	"os"
	"log"
	"io/ioutil"
	"sync"
	"strings"
)

func giveFilesInDir(dir string)([]os.FileInfo, error){
    f, err := os.Open(dir)
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    files, err := f.Readdir(-1)
    f.Close()
    if err != nil {
        log.Fatal(err)
        return nil, err
    }
    return files, nil
}

func printFileNames(files []os.FileInfo){
    for _, file := range files {
        fmt.Println(file.Name())
    }
}

func makeFullPath(src string, files []os.FileInfo)([]string){
	returnArray:= make([]string, len(files))
	for i, file:= range files{
		returnArray[i]= src+"/"+file.Name()
	}
	return returnArray
}

func findInArray(file os.FileInfo, dstFiles []os.FileInfo)(bool){
	for _, element:= range dstFiles{
		if((element.Name()==file.Name())&&(element.ModTime()==file.ModTime())){
			return true
		}
	}
	return false
}

func checkErr(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func copyFile(src string, dst string) {
    data, err := ioutil.ReadFile(src)
    checkErr(err)
    err = ioutil.WriteFile(dst, data, 0644)
    checkErr(err)
}

func filesToCopy(srcFiles []os.FileInfo, dstFiles []os.FileInfo)([]os.FileInfo){
	returnArray:= make([]os.FileInfo,0)
	for _,file:= range srcFiles{
		isPresent:= findInArray(file,dstFiles)
		if(!isPresent){
			returnArray= append(returnArray,file)			
		}
	}
	return returnArray
}

func worker(files []os.FileInfo, srcDir string, dstDir string){
	for _,file:= range files{
		splittedName:= strings.Split(file.Name(),".")
		if(splittedName[len(splittedName)-1]=="ini"){
			continue
		}
		if(file.IsDir()){
			newSrc:= srcDir+"/"+file.Name()
			newDst:= dstDir+"/"+file.Name()
			newContents, err:= giveFilesInDir(newSrc)
			if(err!=nil){
				panic("Directory Reading Error")
				return
			}
			if _,err:= os.Stat(newDst); os.IsNotExist(err){
				os.Mkdir(newDst,0700)
			}
			worker(newContents,newSrc,newDst)
		}else{
			dstFiles, err:= giveFilesInDir(dstDir)
			if(err!=nil){
				panic("Directory Reading Error")
				return
			}
			isPresent:= findInArray(file,dstFiles)
			if(!isPresent){
				copyFile(srcDir+"/"+file.Name(),dstDir+"/"+file.Name())
				fmt.Println(file.Name()+" Copied")
			}
		}
	}
}

func makeDstFolder(srcDir string, dstDir string)(string){
	srcContents:= strings.Split(srcDir,"/")
	srcFolder:= srcContents[len(srcContents)-1]
	newDst:= dstDir+"/"+srcFolder
	if _,err:= os.Stat(newDst); os.IsNotExist(err){
		os.Mkdir(newDst,0700)
	}
	return newDst
}

func handler(allFiles []os.FileInfo, src string, dst string)(bool){
	chanArray:= make([]chan bool,0)
	var mutex = &sync.Mutex{}
	if(len(allFiles)<100){
		for i,file:= range allFiles{
			fmt.Println(i)
			fileArray:= make([]os.FileInfo,1,1)
			fileArray[0]= file
			go func(array []chan bool){
				index:= i
				worker(fileArray,src,dst)
				mutex.Lock()
				channel:= make(chan bool)
				channel<-true
				array= append(array,channel)
				mutex.Unlock()
				fmt.Println("done. new size: ", len(array))
				fmt.Println("pushed at: ", index, "len: ", len(chanArray))
			}(chanArray)
		}
		for(true){
			if(len(allFiles)<=len(chanArray)){
				break
			}
		}
	}else{

	}
	return true
}

func main()(){
	fmt.Println("Hello")
    srcDir:= "C:/Users/usama/Documents/Resume & Applications"
    srcFiles, srcErr:= giveFilesInDir(srcDir)
    dstDir:= "C:/Users/usama/OneDrive"
    if(srcErr!=nil){
    	panic("Directory Reading Error")
    	return
    }
    folderDst:= makeDstFolder(srcDir,dstDir)
    worker(srcFiles,srcDir,folderDst)
}