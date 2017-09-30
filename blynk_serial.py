import getopt
import requests
import serial
import sys
import time

comm = '/dev/ttyACM0'
baud = 9600
token = "your token"
blynk_host = "http://blynk-cloud.com/"
monitor_host = "http://localhost/meteo/"
blynk_disable = False
monitor_disable = False

usage='''
    You can specify port, baud rate, and server endpoint like this:
      python blynk_serial.py -c <serial port> -b <baud rate> -t <your token> \\
      -bs <server address blynk> -ms <server address monitor>

    The defaults are:
      -c,--comm      /dev/ttyUSB0       (on Linux)
                     COM1               (on Windows)
                     /dev/tty.usbserial (on OSX)
      -b,--baud      9600
      -t,--token     "your token"
      -bs,--blynk    blynk-cloud.com
      -ms,--monitor  localhost
      -bdis          blynk disable
      -mdis          monitor disable
'''

try:
    opts,args = getopt.getopt(sys.argv[1:], 'c:b:t:bs:ms:bdis:mdis', 
        ['comm','baud','token','blynk','monitor'])
except:
    print(usage)
    sys.exit(2)

for opt, arg in opts:
    if opt in ('-h', '--help'):
        print(usage)
        sys.exit(2)
    elif opt in ('-c', '--comm'):
        comm = arg
    elif opt in ('-b', '--baud'):
        baud = arg
    elif opt in ('-t', '--token'):
        token = arg
    elif opt in ('-bs', '--blynk'):
        blynk_host = arg
    elif opt in ('-ms', '--monitor'):
        monitor_host = arg
    elif opt in ('-bdis'):
        blynk_disable = True
        print('Blynk disable')
    elif opt in ('-mdis'):
        monitor_disable = True
        print('Monitor disable')

upload_link = token+"/update/"
try:
    ser = serial.Serial(comm)
    ser.baudrate = baud
except Exception:
    exit('Not connect to'+comm)

while True:
    line = ser.readline().decode("utf-8").strip()
    try:
        data = line.split(': ')
        # print(upload_link+data[0]+"?value="+data[1])
        if not blynk_disable:
            try:
                requests.get(blynk_host+upload_link+data[0]+"?value="+data[1])
            except Exception:
                pass

        if not monitor_disable:
            try:
                requests.get(monitor_host+upload_link+data[0]+"?value="+data[1])
            except Exception:
                pass
    except IndexError:
        pass

    print(time.ctime(), line)