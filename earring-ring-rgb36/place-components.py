import math
import pcbnew

pcb = pcbnew.GetBoard()

# Place the LEDs.
radius = 10.5
for n in range(36):
    fp = pcb.FindFootprintByReference('D%d' % (n+1))
    fp.SetPos(pcbnew.VECTOR2I_MM(math.sin(math.pi*2/36*n)*radius, -math.cos(math.pi*2/36*n)*radius))
    fp.SetOrientationDegrees(n*-10-90+360)

# Refresh editor (so it doesn't work on an old version of the PCB).
pcbnew.Refresh()
