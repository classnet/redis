goredis


package main
import (
    "github.com/lxn/walk"
    . "github.com/lxn/walk/declarative"
)
func main() {  
    var label *walk.Label
    Label1:=Label{
      AssignTo: &label,	
      Text:     "Hello world!",			
    }   
    MainForm:=	MainWindow{
        Title:   "windows",
        MinSize: Size{300, 200},          
        Layout:  VBox{},
        Children: []Widget{           
            Label1,
        },		
    }   
    MainForm.Run()
}
