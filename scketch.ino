#include "DHT.h"
#include "MQ135.h"
#include "Wire.h" //for I2C
#include "Adafruit_BMP085.h"
#include "SPI.h" //for SD
#include "SD.h"
#include "SimpleTimer.h"

#define DHTPIN 4 //Датчик температуры и влажности
#define MQPIN A0 //Датчик газа
#define LRPIN A1 //Фоторезистор
#define SDPIN 10 //SDкарта
#define UPDATECO2 60 //Измерение CO2
#define REDLED 7 //Красный цвет
#define GREENLED 6 //Зеленый цвет
#define BLUELED 5 //Синий цвет

String VTEMPERATURE = "V0"; //Температура
String VHUMIDITY = "V1"; //Влажность
String VCO2 = "V2"; //CO2
String VPRESSURE = "V3"; //Давление
String VLIGHT = "V4"; //Освещенность
String VGOOD = "V5"; //Порог хорошегоCO2
String VBAD = "V6"; //Порог плохого CO2
String VBCO2 = "V11"; //Посчитанный C02

DHT dht(DHTPIN, DHT22);
MQ135 gasSensor = MQ135(MQPIN);
Adafruit_BMP085 bmp;
int last_update = -60;
File csvFile;
int iteration = 1;
int failCard = 1;
SimpleTimer timer;
int goodCO2 = 125;
int badCO2 = 150;

void timerEvent(){
  int l = map(analogRead(LRPIN),0,1023, 100, 0);
  float h = dht.readHumidity();
  float t = dht.readTemperature();
  int pressure = bmp.readSealevelPressure()/133.3;
  String strBuff;

  strBuff = VLIGHT + ": " + l;
  Serial.println(strBuff);
  strBuff = VHUMIDITY + ": " + h;
  Serial.println(strBuff);
  strBuff = VTEMPERATURE + ": " + t;
  Serial.println(strBuff);
  strBuff = VPRESSURE + ": " + pressure;
  Serial.println(strBuff);
  
  if((millis()/1000 - last_update) >= UPDATECO2){ //Измерение CO2
    last_update = millis()/1000;
    float co2 = analogRead(MQPIN);

    digitalWrite(REDLED, LOW);
    digitalWrite(GREENLED, LOW);
    digitalWrite(BLUELED, LOW); 

    if(co2 <= goodCO2)
      digitalWrite(GREENLED, HIGH);
    else if(co2 >= badCO2)
      digitalWrite(REDLED, HIGH);
    else
      digitalWrite(BLUELED, HIGH);

    strBuff = VCO2 + ": " + co2;
    Serial.println(strBuff);
    strBuff = VBCO2 + ": " + gasSensor.getCorrectedPPM(t, h);
    Serial.println(strBuff);

    csvFile = SD.open("meteob.csv", FILE_WRITE);
    if(csvFile){
      csvFile.print(iteration);
      csvFile.print("\t");
      csvFile.print(millis()/1000);
      csvFile.print("\t");
      csvFile.print(h);
      csvFile.print("\t");
      csvFile.print(t);
      csvFile.print("\t");
      csvFile.print(l);
      csvFile.print("\t");
      csvFile.print(co2);
      csvFile.print("\t");
      csvFile.print(gasSensor.getCorrectedPPM(t, h));
      csvFile.print("\t");
      csvFile.print(gasSensor.getCorrectedRZero(t, h));
      csvFile.print("\t");
      csvFile.print(pressure);
      csvFile.print("\t");
      csvFile.println();
      csvFile.close();
    } else{
      Serial.println("SD fail");
    }
  }

}

void setup() {
  pinMode(REDLED, OUTPUT);
  pinMode(GREENLED, OUTPUT);
  pinMode(BLUELED, OUTPUT);

  dht.begin();
  pinMode(LRPIN,INPUT);
  Wire.begin();
  bmp.begin();

  pinMode(SDPIN, OUTPUT);
  failCard = !SD.begin(SDPIN);
  if(!failCard){
    csvFile = SD.open("meteob.csv", FILE_WRITE);
    if(csvFile){
      csvFile.print("Iteration\t");
      csvFile.print("Time\t");
      csvFile.print("Humidity\t");
      csvFile.print("Temperature\t");
      csvFile.print("Light\t");
      csvFile.print("analogCO2\t");
      csvFile.print("calculateCO2\t");
      csvFile.print("rZero\t");
      csvFile.print("Pressure\t");
      csvFile.println();
      csvFile.close();
    }
  }

  timer.setInterval(120000L, timerEvent); //2min
  Serial.begin(9600);
  Serial.println("Start");
  timerEvent();
}

void loop() {
  timer.run();
}