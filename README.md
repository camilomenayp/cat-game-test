# cat-game-test
This is just an experiment to practice Go deeper.
It uses the Pixel library ( https://github.com/faiface/pixel ).

Work in progress! so some parts of the code might be really ugly.
## 2020-01-20
Added a better loader of sprites (also using a CSV file). 
Added some basic behaviour when pressing keys, an extremely simple physics struct, and movement.

## 2020-01-16
I separated the code into two packages, one for the main and another to handle the loading of the sprites from the .png files.
Right now, the program is loading a spritesheet with one action ("Walk"), then it shows on screen the animation (animation is only 10 frames, but FPS is way higher, so you see the effect of a "extremely fast cat running". I will work on that next.



### Resources

All sprites are from: 

https://www.gameart2d.com

SDL C++ tutorial (good for theory):

http://lazyfoo.net/tutorials/SDL/i

Article about scrolling and cameras in 2d game:

https://www.gamasutra.com/blogs/ItayKeren/20150511/243083/Scroll_Back_The_Theory_and_P%20ractice_of_Cameras_in_SideScrollers.php
