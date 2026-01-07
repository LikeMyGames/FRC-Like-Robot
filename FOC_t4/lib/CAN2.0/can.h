#ifndef CAN_H_
#define CAN_H_

#include <Arduino.h>
#include <FlexCAN_t4.h>
#include <string>

void initCan();

void addCanRegister(uint8_t id,
                    std::string name,
                    std::function<std::vector<uint8_t>()> readAction,
                    std::function<bool(uint8_t[])> writeAction);
uint32_t getID();
void resetTimeout();
void checkTimeout();
void canSniff(const CAN_message_t &);
void ProcessIdSpecificMessage(const CAN_message_t &);
void canLoop();
void Status();

std::vector<uint8_t> CanEchoRead();
bool CanEchoWrite(uint8_t args[]);

#endif