# gograceful
Module to handle graceful shutdown of golang applications. 

### Install
```
go get github.com/dreson4/gograceful
```

### How to use
In your main, call this function, it will intercept signal terminate syscalls and prevent shutdown until all operations have finished.
The ``func(){}`` passed will be called when finally shutting down.
```
package main

func init(){
  gograceful.HandleGracefulShutdown(func(){})
}
```

```
  func importantOperation(){
    shouldRun := gograceful.AddRunningOperation()
    if !shouldRun{
      //in the middle of shutting down, no new operation is accepted.
      return
    }
    defer gograceful.FinishRunningOperation()
    
    //Run anything here, won't be terminated until the function returns. 
  }

```
