import sys
import PIL
from PIL import Image
import os

def filter_hidden(s):
  if s.startswith('.'):
    return False
  else:
    return True

name = 'blocks.png'
path = "/Users/wegesdal/Documents/25-ground-blocks/"
blocks = sorted([q for q in filter(filter_hidden, [p for p in os.listdir(path)])])

print(blocks)

target_size = 128

new_im = Image.new('RGBA', (target_size*5, target_size*5))

max_width = 0
max_height = 0
for i in range(len(blocks)):
    current_block = Image.open(path+blocks[i])
    if current_block.size[0] > max_width:
        max_width = current_block.size[0]
    if current_block.size[1]>max_height:
        max_height = current_block.size[1]

ratio = min(target_size / max_width, target_size / max_height)
print(ratio)
for i in range(len(blocks)):
    current_block = Image.open(path+blocks[i])
    current_block = current_block.resize((int(current_block.size[0]*ratio), int(current_block.size[1]*ratio)), resample=PIL.Image.LANCZOS)
    new_im.paste(current_block, ((i % 5)*target_size, (i // 5)*target_size))

new_im.save('assets/sprites/'+name)