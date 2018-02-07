package main

import(
  "io/ioutil"
  "os"
  "testing"
  "path"
  "strings"
)

const (
  tempIOPrefix = "common_"
  tempFileName = "common"
  tempFileNameLCSF = "common.lcsf"
  encryptedFileExtension = "lcsf"
  genNumOfBytes = 100
)

const (
  createTempFileTempDirError = "Unable to create temporary file/directory"
  closeTempFileError = "Unable to close temporary file"
  removeTempFileTempDirError = "Unable to remove temporary file/directory"
  createTempLCSFFile = "Unable to create temporary file with LCSF extension"
  closeTempLCSFFile = "Unable to close temporary file with LCSF extension"
  removeTempLCSFFile = "Unable to remove temporary file LCSF extension"
)

// HandleError will handle error returned from OS calls that shouldn't
// return nil
func HandleError(t *testing.T, error string, err error) {
  t.Helper()
  if err != nil {
    t.Errorf("%s", error)
    t.FailNow()
  }
}

// TestFileExits tests two cases, the return value of fileExists when passing it a
// known existing file, and when passing it a file we know does/should not exist
func TestFileExists (t *testing.T) {
  file, err := ioutil.TempFile("", tempIOPrefix)
  HandleError(t, createTempFileTempDirError, err)

  t.Logf("%s was created!", file.Name())

  if exists, _ := fileExists(file.Name()); exists != true {
    t.Error("FAILED: fileExistsWithValidFileReturnsFalse")
  }

  // This file should not exist
  if exists, _ := fileExists(file.Name() + "meow"); exists != false {
    t.Error("FAILED: fileExistsWithInvalidFileReturnsTrue")
  }

  HandleError(t, closeTempFileError, file.Close())
  HandleError(t, removeTempFileTempDirError, os.Remove(file.Name()))

  t.Logf("%s was removed!", file.Name())

  if exists, _ := fileExists(""); exists != false {
    t.Error("FAILED: fileExistsEmptyFileNameReturnsTrue")
  }
  //Max file length test?
}

// TestGenerateRandomBytes tests two cases, that generateRandomBytes returns the same amount of bytes
// that is passed as the length and that not all bytes in the returns array is 0.
func TestGenerateRandomBytes(t *testing.T) {
  genBytes, _ := generateRandomBytes(genNumOfBytes)
  if len(genBytes) != genNumOfBytes {
    t.Error("FAILED: generateRandomBytesIncorrectLength")
  }

  for i := range genBytes {
    if genBytes[i] != 0 {
      break
    } else if i + 1 == genNumOfBytes {
      t.Error("FAILED: generateRandomBytesGeneratesZeros")
    }
  }
  // Can we force to return error?
}


func TestGetDecryptedFilename(t *testing.T) {
  // Test empty string cases
  file, err := getDecryptedFilename("", "")

  // Test removing extension from filename without providing an output path
  file, err = getDecryptedFilename(tempFileNameLCSF, "")

  if strings.HasSuffix(file, encryptedFileExtension)  {
    // This fails if you encrypt a lcsf file...
    t.Error("FAILED: getDecryptedFilenameDoesntRemoveLCSFWithoutOutputPath")
  }

  tempDir, err := ioutil.TempDir("", tempIOPrefix)
  HandleError(t, createTempFileTempDirError, err)

  // Test when path doesn't end with lcsf
  // Creates "common" file in tempDir
  fullFileName := path.Join(tempDir, tempFileName)
  f, err := os.OpenFile(fullFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_SYNC, 0600)
  HandleError(t, createTempFileTempDirError , err)
  HandleError(t, closeTempFileError , f.Close())

  file, err = getDecryptedFilename(fullFileName, tempDir)

  if err == nil {
    t.Error("FAILED: getDecryptedFilenameReturnsNoErrorForFileWithoutLCSF")
  }

  HandleError(t, removeTempFileTempDirError, os.Remove(fullFileName))

  // Test removing extension from filename, providing output path
  // Creates "common.lcsf" file in tempDir
  fullFileName = path.Join(tempDir, tempFileNameLCSF)
  f, err = os.OpenFile(fullFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_SYNC, 0600)
  HandleError(t, createTempLCSFFile , err)
  HandleError(t, closeTempLCSFFile , f.Close())

  file, err = getDecryptedFilename(fullFileName, tempDir)
  HandleError(t, createTempLCSFFile , err)

  if strings.HasSuffix(file, encryptedFileExtension)  {
    // This fails if you encrypt a lcsf file...
    t.Error("FAILED: getDecryptedFilenameDoesntRemoveLCSF")
  }

  HandleError(t, removeTempLCSFFile, os.Remove(fullFileName))

  // Test removing extension from provided a regular file output path?
  // Creates "common.lcsf" file in tempDir
  fullFileName = path.Join(tempDir, tempFileNameLCSF)
  f, err = os.OpenFile(fullFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_SYNC, 0600)
  HandleError(t, createTempLCSFFile , err)
  HandleError(t, closeTempLCSFFile , f.Close())

  // Not sure what the use case is for this, but it hits the IsReg condition
  file, err = getDecryptedFilename(fullFileName, fullFileName)
  HandleError(t, createTempLCSFFile , err)

  if strings.HasSuffix(file, encryptedFileExtension) {
    // This fails if you encrypt a lcsf file...
    t.Error("FAILED: getDecryptedFilenameDoesntRemoveExtensionFromRegularOutpath")
  }
  err = os.Remove(tempDir)
  // Remove test directory
  if err != nil {
    t.Error("Could not clean up temporary directory, manually clean and fix")
    HandleError(t, removeTempFileTempDirError, err)
  }
}

func TestGetEncryptedFilename(t *testing.T) {
  // Test invalid file name
  // TODO: Should be called/completed when returning a proper error
  //_ = getEncryptedFilename("", "")

  // Test valid file name
  file, err := ioutil.TempFile("", tempIOPrefix)
  HandleError(t, createTempFileTempDirError, err)

  t.Logf("%s was created!", file.Name())

  f := getEncryptedFilename(file.Name(), "")
  if !strings.HasSuffix(f, encryptedFileExtension) {
    // This fails if you encrypt a lcsf file...
    t.Error("FAILED: getEncryptedFilenameDoesntAppendExtensionFromRegularOutpath")
  }

  // Test valid file name with already existing .lcsf file extension
  // TODO: Define if this should be okay or not
  //temp := getEncryptedFilename(f, "")

  HandleError(t, closeTempFileError, file.Close())
  HandleError(t, removeTempFileTempDirError, os.Remove(file.Name()))

  t.Logf("%s was removed!", file.Name())
  // Test valid file; valid output path, is a directory

  // Test valid file; valid output path, is not a directory

  // Test valid file; invalid output path-doesn't exist, no output path extension

  // Test valid file; invalid output path-doesn't exist, output path extension exists

  // Test invalid output path--exists but not directory, valid file


}
