import sys
from PIL import Image
import os

def filter_hidden(s):
  if s.startswith('.'):
    return False
  else:
    return True

name = 'doodads.png'
path = "/Users/wegesdal/Documents/blender/doodads/render/"
doodads = [q for q in filter(filter_hidden, [p for p in os.listdir(path)])]

print(doodads)

sprite_size = 512
max_height = sprite_size * len(doodads)
number_of_angles = 8
number_of_frames_per_angle = 1

new_im = Image.new('RGBA', (sprite_size*number_of_angles*number_of_frames_per_angle, max_height))
for i in range(len(doodads)):
  y_offset = i * sprite_size
  angles = [r for r in filter(filter_hidden, [s for s in os.listdir(path + doodads[i])])]
  print(angles)
  for a in range(len(angles)):
    pics_list = sorted(filter(filter_hidden, [p for p in os.listdir(path + doodads[i] + '/' + angles[a])]))
    pics = [Image.open(path+doodads[i]+'/'+angles[a]+'/'+p) for p in pics_list]
    max_width = sprite_size * len(pics) * len(angles)
    for p in range(len(pics)):
      # paste
      x_offset = a*len(pics)*sprite_size + p*sprite_size
      new_im.paste(pics[p], (x_offset,y_offset))

new_im.save('assets/sprites/'+name)