package main

import (
    "bufio"
    "fmt"
    "syscall"
    "os/exec"
    "os"
    "bytes"
    "time"
    "strings"
    "log"
)

func getInfo() (string, string){

  scanner := bufio.NewScanner(os.Stdin)

  fmt.Print("Enter location of lol.launcher.exe: ")
  scanner.Scan()
  launcherLocale := scanner.Text()
  launcherLocale = launcherLocale + "\\lol.launcher.exe"

  fmt.Print("Enter Summoner Name: ")
  scanner.Scan()
  sumName := scanner.Text()

  return launcherLocale, sumName
}

func getLink(sumName string) (string){

  sumName = strings.Replace(sumName, " ","+", -1)
  link := "http://www.lolnexus.com/NA/search?name=" + sumName + "&region=NA"

  return link
}

func isProcRunning(names ...string) (bool, error) {
    if len(names) == 0 {
        return false, nil
    }

    cmd := exec.Command("tasklist.exe", "/fo", "csv", "/nh")
    cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
    out, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Println(out)
        return false, err
    }

    for _, name := range names {
        if bytes.Contains(out, []byte(name)) {
            return true, nil
        }
    }
    return false, nil
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func readLines(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines, scanner.Err()
}

func write(lines []string, path string) error {
  file, err := os.Create(path)
  if err != nil {
    return err
  }
  defer file.Close()

  w := bufio.NewWriter(file)
  for _, line := range lines {
    fmt.Fprintln(w, line)
  }
  return w.Flush()
}

func startClient(){
  clientRunning, err2 := isProcRunning("LolClient.exe","lolpatcher.exe","LoLPatcherUx.exe")

  if !clientRunning{
    c := exec.Command("cmd", "/C", launcherLocale)
    if err := c.Run(); err != nil {
        fmt.Println("Error: ", err, err2)
    }
  }
}

func main(){

  var launcherLocale, sumName string
  empty := true

  if _, err := os.Stat("config.txt"); err == nil {
    lines, err := readLines("config.txt")

    if err != nil {
      log.Fatalf("readLines: %s", err)
    }

    if lines[1] != ""{
      launcherLocale = lines[0]
      sumName = lines[1]
      empty = false
    }
  }

  if empty {
    launcherLocale, sumName = getInfo()
    info := []string{launcherLocale, sumName}
    write(info, "config.txt")
  }

  startClient()

  t := time.NewTicker(10 * time.Second)
  ran := false
  for now := range t.C {
    isRunning, err := isProcRunning("League of Legends.exe")

    if isRunning && !ran{
        link := getLink(sumName)
        exec.Command("cmd", "/c", "start", link).Start()
        ran = true
    }
    if !isRunning {ran = false}

    clientRunning, err2 := isProcRunning("LolClient.exe","lolpatcher.exe","LoLPatcherUx.exe")
    if !clientRunning{os.Exit(0)}

    if err != nil || err2 != nil {
      fmt.Println("Error: ", err)
      fmt.Println("Error: ", err2)
      fmt.Println("Time of error ", now)
    }
  }
}
