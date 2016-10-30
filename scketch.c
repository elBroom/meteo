#include "DHT.h"
#include "MQ135.h"
#include "Wire.h" //for I2C
#include "Adafruit_BMP085.h"
#include "SPI.h" //for SD
#include "SD.h"
#include "BlynkSimpleStream.h"
#include "SimpleTimer.h"

#define DHTPIN 4 //Датчик температуры и влажности
#define MQPIN A0 //Датчик газа
#define LRPIN A1 //Фоторезистор
#define SDPIN 10 //SDкарта
#define UPDATECO2 60 //Измерение CO2
#define REDLED 7 //Красный цвет
#define GREENLED 6 //Зеленый цвет
#define BLUELED 5 //Синий цвет

#define VTEMPERATURE V0 //Температура
#define VHUMIDITY V1 //Влажность
#define VCO2 V2 //CO2
#define VPRESSURE V3 //Давление
#define VLIGHT V4 //Освещенность
#define VGOOD V5 //Порог хорошегоCO2
#define VBAD V6 //Порог плохого CO2
#define VBCO2 V11 //Посчитанный C02
#define VLED V10 //Маяк SD записи

#define BTOKEN "your token"

WidgetLED led1(VLED);
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

  Blynk.virtualWrite(VLIGHT, l);
  Blynk.virtualWrite(VHUMIDITY, h);
  Blynk.virtualWrite(VTEMPERATURE, t);
  Blynk.virtualWrite(VPRESSURE, pressure);
  
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

    Blynk.virtualWrite(VCO2, co2);
    Blynk.virtualWrite(VBCO2, gasSensor.getCorrectedPPM(t, h));

    csvFile = SD.open("meteob.csv", FILE_WRITE);
    if(csvFile){
      led1.on();
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
      led1.off();
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

  Serial.begin(9600);
  Blynk.begin(BTOKEN, Serial);

  led1.off();

  pinMode(SDPIN, OUTPUT);
  failCard = !SD.begin(SDPIN);
  if(!failCard){
    csvFile = SD.open("meteob.csv", FILE_WRITE);
    if(csvFile){
      led1.on();
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

  timer.setInterval(20000L, timerEvent); //20sec

  Blynk.run();
  Blynk.virtualWrite(VGOOD, goodCO2);
  Blynk.virtualWrite(VBAD, badCO2);
}

BLYNK_WRITE(VGOOD){
  goodCO2 = param.asInt();
}

BLYNK_WRITE(VBAD){
  badCO2 = param.asInt();
}

void loop() {
  Blynk.run();
  timer.run();
}