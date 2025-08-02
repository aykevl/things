#!/usr/bin/env python3

import math

from kipy import KiCad
from kipy.util import from_mm
from kipy.geometry import Angle, Vector2
from kipy.board_types import BoardCircle

BL_F_SilkS = 40

dot_radius = 300000
dot_multiply = 2200000
pcb_radius = 20e6

kicad = KiCad()
board = kicad.get_board()

# Determine where the LEDs should go
fibonacci_positions = []
start = 4
end = start + 64
rotate = 55
for i in range(start, end):
    goldenAngle = 180 * (3 - math.sqrt(5))
    r = math.sqrt(i)
    angle = (i * goldenAngle + rotate) % 360 
    x = r * math.cos(angle / 180 * math.pi);
    y = -r * math.sin(angle / 180 * math.pi);
    fibonacci_positions.append((int(x*dot_multiply), int(y*dot_multiply)))

# Remove all existing dots.
def replace_dots():
    dots = []
    for item in board.get_shapes():
        if item.layer == BL_F_SilkS and item.radius() == dot_radius:
            dots.append(item)
    if len(dots):
        board.remove_items(dots)

    # Create all dots.
    items = []
    for (x, y) in fibonacci_positions:
        dot = BoardCircle()
        dot.center = Vector2.from_xy(x, y)
        dot.radius_point = Vector2.from_xy(dot.center.x, dot.center.y+dot_radius)
        dot.layer = BL_F_SilkS
        dot.attributes.fill.filled = True
        items.append(dot)
    board.create_items(items)

def place_leds():
    update = []
    for fp in board.get_footprints():
        # Filter for LEDs inside the PCB.
        if not fp.reference_field.text.value.startswith('D'):
            continue
        if math.fabs(fp.position.length()) > pcb_radius:
            continue

        # determine closest dot
        distance = 1e9 # 1m
        closest = None
        for (x, y) in fibonacci_positions:
            dx = math.fabs(x-fp.position.x)
            dy = math.fabs(y-fp.position.y)
            d = math.sqrt(dx*dx + dy*dy)
            if d < distance:
                distance = d
                closest = (x, y)

        position = Vector2.from_xy(closest[0], closest[1])
        print('position:', position.angle_degrees(), fp.orientation.normalize().degrees)
        degrees = -position.angle_degrees()
        old_degrees = fp.orientation.normalize().degrees
        if degrees % 90 == old_degrees % 90:
            degrees = old_degrees

        if fp.position.x != closest[0] or fp.position.y != closest[1] or old_degrees != degrees:
            print('update:', fp.reference_field.text.value)
            fp.position = position
            fp.orientation = Angle.from_degrees(degrees)
            update.append(fp)

    if len(update):
        board.update_items(update)
        print('updated')

# Reset LED positions, put them organized outside the PCB.
def organize_leds():
    update = []
    for fp in board.get_footprints():
        # Filter for LEDs inside the PCB.
        reference = fp.reference_field.text.value
        if not reference.startswith('D'):
            continue
        n = int(reference[1:])
        row = (n-1) % 16
        column = (n-1) // 16
        fp.orientation = Angle.from_degrees(0)
        fp.position = Vector2.from_xy(int((row - 8) * 2e6), int(column * 2e6 - 28e6))
        update.append(fp)
        print(n, row, column)
    board.update_items(update)

if __name__ == '__main__':
    #organize_leds()
    place_leds()
