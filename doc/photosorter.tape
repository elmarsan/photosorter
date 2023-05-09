Output ./doc/photosorter.gif
Set Height 1200
Set Width 1200
Set Margin 60
Set MarginFill "./doc/blue.jpg"
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

Type "tree ./test/img" 
Sleep 1.5s 
Enter
Sleep 3s 

Type "photosorter sort ./test/img ./sorted_imgs --format=month" 
Sleep 3s
Enter
Sleep 3s

Type "tree ./sorted_imgs" 
Sleep 3s
Enter
Sleep 5s

Hide 
Type "rm -r ./sorted_imgs"
Enter