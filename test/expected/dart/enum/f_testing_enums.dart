// Autogenerated by Frugal Compiler (2.2.2)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

enum testing_enums {
  one,
  two,
  Three,
}

int serializetesting_enums(testing_enums variant) {
  switch (variant) {
    case testing_enums.one:
      return 45;
    case testing_enums.two:
      return 3;
    case testing_enums.Three:
      return 76;
  }
}

testing_enums deserializetesting_enums(int value) {
  switch (value) {
    case 45:
      return testing_enums.one;
    case 3:
      return testing_enums.two;
    case 76:
      return testing_enums.Three;
    default:
      throw new thrift.TProtocolError(thrift.TProtocolErrorType.UNKNOWN, "Invalid value '$value' for enum 'testing_enums'");  }
}
