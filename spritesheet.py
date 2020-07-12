import sys
from PIL import Image
import os

def filter_hidden(s):
  if s.startswith('.'):
    return False
  else:
    return True

path = "/Users/wegesdal/Documents/blender/terminal/render/"
poses = [q for q in filter(filter_hidden, [p for p in os.listdir(path)])]

print(poses)

sprite_size = 128
max_height = sprite_size * len(poses)
number_of_angles = 4
number_of_frames_per_angle = 10

new_im = Image.new('RGBA', (sprite_size*number_of_angles*number_of_frames_per_angle, max_height))
for i in range(len(poses)):
  y_offset = i * sprite_size
  angles = [r for r in filter(filter_hidden, [s for s in os.listdir(path + poses[i])])]
  print(angles)
  for a in range(len(angles)):
    pics_list = sorted(filter(filter_hidden, [p for p in os.listdir(path + poses[i] + '/' + angles[a])]))
    pics = [Image.open(path+poses[i]+'/'+angles[a]+'/'+p) for p in pics_list]
    max_width = sprite_size * len(pics) * len(angles)
    for p in range(len(pics)):
      # paste
      x_offset = a*len(pics)*sprite_size + p*sprite_size
      new_im.paste(pics[p], (x_offset,y_offset))

new_im.save('terminal.png')