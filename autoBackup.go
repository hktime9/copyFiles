package main

import(
	"fmt"
	"os"
	"log"
	"io/ioutil"
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

// func handler(paths []string)(int){
// 	if(len(paths)<100){
		
// 	}else{
		
// 	}
// }

func checkErr(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func Copy(src string, dst string) {
    // Read all content of src to data
    data, err := ioutil.ReadFile(src)
    checkErr(err)
    // Write data to dst
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
				Copy(srcDir+"/"+file.Name(),dstDir+"/"+file.Name())
				fmt.Println(file.Name()+" Copied")
			}
		}
	}
}

func main()(){
	fmt.Println("Hello")
    srcDir:= "C:/Users/usama/Code/Go/AutoBackup/src"
    srcFiles, srcErr:= giveFilesInDir(srcDir)
    dstDir:= "C:/Users/usama/Code/Go/AutoBackup/dst"
    dstFiles, dstErr:= giveFilesInDir(dstDir)
    if(srcErr!=nil||dstErr!=nil){
    	panic("Directory Reading Error")
    	return
    }
    files:= filesToCopy(srcFiles,dstFiles)

    // newStr:= "C:/Users/usama/Code/Go/AutoBackup/src/9701_s07_qp_5.pdf"
    worker(files,srcDir,dstDir)
    // absolutePaths:= makeFullPath(srcDir,files)
}