Output ./photosorter.gif
Set Height 1200
Set Width 1200
Set Margin 60
Set MarginFill "blue.jpg"
Set WindowBar "Rings"

Set WindowBar Colorful
Set WindowBarSize 40
Set BorderRadius 10

Set Padding 20

Require tree

Type "go install github.com/elmarsan/photosorter@v1.0.0" 
Sleep 2s 
Enter
Sleep 1s

Type "tree ./img" 
Sleep 1.5s 
Enter
Sleep 3s 

Type "photosorter sort ./img ./sorted_imgs --format=month" 
Sleep 3s
Enter
Sleep 3s

Type "tree ./sorted_imgs" 
Sleep 3s
Enter
Sleep 5s