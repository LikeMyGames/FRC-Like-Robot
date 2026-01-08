#include <can.h>
#include <FlexCAN_T4.h>
#include <EEPROM.h>
#include <TeensyThreads.h>
#include <util_address.h>
#include <motor.h>
#include <string>

#define NUM_TX_MAILBOXES 2
#define NUM_RX_MAILBOXES 6

typedef struct
{
    std::string name;
    std::function<std::vector<uint8_t>()> read;
    std::function<bool(uint8_t[])> write;
} CAN_REGISTER_t;

FlexCAN_T4<CAN2, RX_SIZE_256, TX_SIZE_16> Can0;
std::unordered_map<uint8_t, CAN_REGISTER_t *> REG_MAP;
uint32_t CAN_ID = 4;
int LOOP_THREAD_ID;
int TIMEOUT_THREAD_ID;
volatile uint16_t TIMEOUT_MS_COUNT = 0;
volatile static uint16_t MAX_TIMEOUT_MS = 2000;

/// @brief
void initCan()
{
    CAN_ID = EEPROM.read(CAN_ID_ADDRESS);
    MAX_TIMEOUT_MS = EEPROM.read(CAN_MAX_TIMEOUT_MS_ADDRESS);
    delay(400);
    Can0.begin();
    Can0.setBaudRate(1 * 1000 * 1000);
    Can0.setMaxMB(NUM_TX_MAILBOXES + NUM_RX_MAILBOXES);
    for (int i = 0; i < NUM_RX_MAILBOXES; i++)
    {
        Can0.setMB((FLEXCAN_MAILBOX)i, RX, EXT);
    }
    for (int i = NUM_RX_MAILBOXES; i < (NUM_TX_MAILBOXES + NUM_RX_MAILBOXES); i++)
    {
        Can0.setMB((FLEXCAN_MAILBOX)i, TX, EXT);
    }
    Can0.setMBFilter(REJECT_ALL);
    Can0.enableMBInterrupts();
    Can0.onReceive(MB0, canSniff);
    Can0.onReceive(MB3, ProcessIdSpecificMessage);
    Can0.setMBUserFilter(MB3, CAN_ID, 0xff);
    Can0.setMBUserFilter(MB0, 0x00, 0x00);
    Can0.enhanceFilter(MB0);
    Can0.enhanceFilter(MB3);
    Can0.mailboxStatus();
    Can0.distribute();
    Serial.println(getID());
    LOOP_THREAD_ID = threads.addThread(canLoop);
    TIMEOUT_THREAD_ID = threads.addThread(checkTimeout);
    Serial.println(LOOP_THREAD_ID);
    Serial.println(TIMEOUT_THREAD_ID);
}

std::vector<uint8_t> CanEchoRead()
{
    return (std::vector<uint8_t>){0};
}

bool CanEchoWrite(uint8_t args[])
{
    return true;
}

/// @brief
/// @param msg
void canSniff(const CAN_message_t &msg)
{
    digitalToggle(LED_BUILTIN);
    Serial.print("MB ");
    Serial.print(msg.mb);
    Serial.print("  OVERRUN: ");
    Serial.print(msg.flags.overrun);
    Serial.print("  LEN: ");
    Serial.print(msg.len);
    Serial.print(" EXT: ");
    Serial.print(msg.flags.extended);
    Serial.print(" TS: ");
    Serial.print(msg.timestamp);
    Serial.print(" ID: ");
    Serial.print(msg.id, HEX);
    Serial.print(" Buffer: ");
    for (uint8_t i = 0; i < msg.len; i++)
    {
        Serial.print(msg.buf[i], HEX);
        Serial.print(" ");
    }
    Serial.println();
}

void ProcessIdSpecificMessage(const CAN_message_t &msg)
{
    resetTimeout();
    Serial.println("Got ID specific messasge");

    if (msg.len == 0)
    {
        return;
    }

    CAN_message_t rmsg;
    rmsg.id = msg.id;
    rmsg.buf[0] = msg.buf[0];
    rmsg.flags.extended = true;

    if (msg.len == 1)
    {

        if (REG_MAP.count(msg.buf[0]) == 0)
        {
            Serial.print("No register exists with id ");
            Serial.println(msg.buf[0]);
        }

        auto reg = REG_MAP[msg.buf[0]];

        std::vector<uint8_t> buf = reg->read();
        if (buf.size() > 7)
        {
            uint size = buf.size();
            size %= 7;
            uint bytesRead = 0;
            rmsg.len = 8;
            while (buf.size() - bytesRead > size)
            {
                for (int i = 0; i < 7; i++)
                {
                    rmsg.buf[i + 1] = buf.at(bytesRead + i);
                    bytesRead++;
                }

                Can0.write(rmsg);
            }

            if (size > 0)
            {
                rmsg.len = size + 1;
                for (uint i = 0; i < size; i++)
                {
                    rmsg.buf[i + 1] = buf.at(bytesRead + i);
                }
                Can0.write(rmsg);
            }
        }
        else
        {
            rmsg.len = buf.size() + 1;
            for (uint i = 1; i < buf.size(); i++)
            {
                rmsg.buf[i] = buf.at(i - 1);
            }
            Can0.write(rmsg);
        }
        return;
    }

    rmsg.len = 2;
    uint8_t funcBuf[7];
    for (uint i = 0; i < sizeof(funcBuf); i++)
    {
        funcBuf[i] = msg.buf[i + 1];
    }
    auto reg = REG_MAP[msg.buf[0]];
    rmsg.buf[1] = reg->write(funcBuf);

    Can0.write(rmsg);
}

void addCanRegister(uint8_t id, std::string name, std::function<std::vector<uint8_t>()> readAction, std::function<bool(uint8_t[])> writeAction)
{
    if (REG_MAP.count(id) != 0)
    {
        Serial.println("Could not add register, already exists");
    }

    CAN_REGISTER_t *reg = {};
    reg->name = name;
    reg->read = readAction;
    reg->write = writeAction;

    REG_MAP.insert_or_assign(id, reg);
}

void canLoop()
{
    while (1)
    {
        Can0.events();
        if (Can0.getRXQueueCount() != 0)
        {
            Serial.println(Can0.getRXQueueCount());
        }
    }
}

void Status()
{
    Can0.mailboxStatus();
}

uint32_t getID()
{
    return CAN_ID;
}

void setID(uint32_t newID)
{
    Can0.FLEXCAN_EnterFreezeMode();
    Can0.setMBUserFilter(MB3, newID, 0xff);
    Can0.FLEXCAN_ExitFreezeMode();
    CAN_ID = newID;
}

void resetTimeout()
{
    TIMEOUT_MS_COUNT = 0;
}

void checkTimeout()
{
    while (1)
    {
        if (TIMEOUT_MS_COUNT >= MAX_TIMEOUT_MS)
        {
            disableAllMotors();
            threads.delay(500);
        }
        TIMEOUT_MS_COUNT++;
        threads.delay(1);
    }
}

uint16_t getMaxCanTimeout()
{
    return TIMEOUT_MS_COUNT;
}

void setMaxCanTimeout(uint16_t newTimeout)
{
    MAX_TIMEOUT_MS = newTimeout;
}