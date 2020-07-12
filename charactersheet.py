import sys
from PIL import Image
import os

def filter_hidden(s):
  if s.startswith('.'):
    return False
  else:
    return True

path = "/Users/wegesdal/Documents/blender/gopher/render/"
poses = [q for q in filter(filter_hidden, [p for p in os.listdir(path)])]

print(poses)

sprite_size = 256
max_height = sprite_size * len(poses) * 4
print("max_height: {}".format(max_height))
number_of_angles = 8
number_of_frames_per_angle = 10

new_im = Image.new('RGBA', (sprite_size*number_of_angles*number_of_frames_per_angle//4, max_height))
for i in range(len(poses)):
  
  angles = sorted([r for r in filter(filter_hidden, [s for s in os.listdir(path + poses[i])])])
  print(angles)
  for a in range(len(angles)):
    pics_list = sorted(filter(filter_hidden, [p for p in os.listdir(path + poses[i] + '/' + angles[a])]))
    pics = [Image.open(path+poses[i]+'/'+angles[a]+'/'+p) for p in pics_list]
    max_width = sprite_size * len(pics) * len(angles) / 4
    print(max_width)
    num_frames = len(pics)
    for f in range(num_frames):
      # paste
      y_offset = i * sprite_size * 4 + (a // 2) * sprite_size
      x_offset = (a%2)*num_frames*sprite_size + f*sprite_size
      new_im.paste(pics[f], (x_offset,y_offset))

new_im.save('gopher.png')