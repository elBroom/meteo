#include "DHT.h"
//#include "MQ135.h"
#include "Wire.h" //for I2C
#include "Adafruit_BMP085.h"
#include "SPI.h" //for SD
#include "SD.h"

#define DHTPIN 2 //Датчик температуры и влажности
#define MQPIN A0 //Датчик газа
#define LRPIN A1 //Фоторезистор
#define SDPIN 10 //SDкарта

#define REDLED 7 //Красный цвет
#define GREENLED 6 //Зеленый цвет
#define BLUELED 5 //Синий цвет

DHT dht(DHTPIN, DHT22);
//MQ135 gasSensor = MQ135(MQPIN);
Adafruit_BMP085 bmp;
File csvFile;
int iteration = 1;

void setup() {
  pinMode(REDLED, OUTPUT);
  pinMode(GREENLED, OUTPUT);
  pinMode(BLUELED, OUTPUT);
  
  Serial.begin(9600);    
  
  dht.begin();  

  pinMode(LRPIN,INPUT);

  Wire.begin(); 
  bmp.begin();
  
  pinMode(SDPIN, OUTPUT);
  if(!SD.begin(SDPIN)){
    Serial.println("Card Failed");
    return;
  }
  Serial.println("Card Redy");

  csvFile = SD.open("meteo.csv", FILE_WRITE);
  if(csvFile){
    csvFile.print("Iteration\t");
    csvFile.print("Humidity\t");
    csvFile.print("Temperature\t");
    csvFile.print("Light\t");
    csvFile.print("CO2\t");
    csvFile.print("Pressure\t");
    csvFile.print("Temperature(bmp)\t");
    csvFile.print("Altitude\t");
    csvFile.println();
    csvFile.close();
  }
}

void loop() {
  float sealevelPressure = 100600.00;
  float goodCO2 = 120.0;
  float badCO2 = 150.0;
  
  int l = analogRead(LRPIN);
  float h = dht.readHumidity();
  float t = dht.readTemperature();
  float co2 = analogRead(MQPIN);
  float pressure = bmp.readSealevelPressure()/133.3;
  float temperature = bmp.readTemperature();
  float altitude = bmp.readAltitude(sealevelPressure);

  digitalWrite(REDLED, LOW);
  digitalWrite(GREENLED, LOW);
  digitalWrite(BLUELED, LOW);

  if(co2 <= goodCO2)
    digitalWrite(GREENLED, HIGH);
  else if(co2 >= badCO2)
    digitalWrite(REDLED, HIGH);
  else
    digitalWrite(BLUELED, HIGH);

  Serial.print("Iteration: ");
  Serial.print(iteration);
  Serial.print("\t");
  Serial.print("Time: ");
  Serial.print(millis()/1000);
  Serial.print("\t");
  Serial.print("Humidity: ");
  Serial.print(h);
  Serial.print("%\t");
  Serial.print("Temperature: ");
  Serial.print(t);
  Serial.print("*C\t");
  Serial.print("Light: ");
  Serial.print(map(l,0,1023, 100, 0));
  Serial.print("%\t");
  Serial.print("CO2: ");
  Serial.print(co2);
  Serial.print("\t");
  Serial.print("Pressure: ");
  Serial.print(pressure);
  Serial.print("mm Hg\t");
  Serial.print("Temperature(bmp): ");
  Serial.print(temperature);
  Serial.print("*C\t");
  Serial.print("Altitude: ");
  Serial.print(altitude);
  Serial.print("m\t");
  Serial.println();

  csvFile = SD.open("meteo.csv", FILE_WRITE);
  if(csvFile){
    csvFile.print(iteration);
    csvFile.print("\t");
    csvFile.print(millis()/1000);
    csvFile.print("\t");
    csvFile.print(h);
    csvFile.print("\t");
    csvFile.print(t);
    csvFile.print("\t");
    csvFile.print(map(l,0,1023, 100, 0));
    csvFile.print("\t");
    csvFile.print(co2);
    csvFile.print("\t");
    csvFile.print(pressure);
    csvFile.print("\t");
    csvFile.print(temperature);
    csvFile.print("\t");
    csvFile.print(altitude);
    csvFile.print("\t");
    csvFile.println();
    csvFile.close();
  }

  iteration++;
  delay(180000); //3 min
}