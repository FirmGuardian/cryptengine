package main

import(
  "testing"
  "io/ioutil"
  "os"
  "path"
)

func CreateZipFile(t *testing.T) string {
  t.Helper()
  var files []string
  file, err := ioutil.TempFile("", tempIOPrefix)
  HandleError(t, createTempFileTempDirError, err)
  HandleError(t, closeTempFileError, file.Close())

  t.Logf("%s was created!", file.Name())

  files = append(files, file.Name())
  return archiveFiles(files)
}

func TestArchiveFiles (t *testing.T) {
  var files []string
  // ---------------------------------------
  // Pass existing files
  // ---------------------------------------
  file, err := ioutil.TempFile("", tempIOPrefix)
  HandleError(t, createTempFileTempDirError, err)
  HandleError(t, closeTempFileError, file.Close())

  t.Logf("%s was created!", file.Name())

  files = append(files, file.Name())
  archivePath1 := archiveFiles(files)
  if archivePath1 == "" {
    t.Error("FAILED: archiveFilesCantZipValidFile")
  } else {
    t.Logf("%s was created!", archivePath1)
  }

  // ---------------------------------------
  // Pass non-existing files
  // ---------------------------------------
  files = append(files, "1")
  // Drop the first element (an existing path)
  failedArchivePath := archiveFiles(files[1:])
  if failedArchivePath != "" {
    t.Error("FAILED: archiveFilesZipsNonexistantFiles")
  }

  // ---------------------------------------
  // Pass a mix of existing and non-existing
  // ---------------------------------------
  archivePath2 := archiveFiles(files)
  if archivePath2 == "" {
    t.Error("FAILED: archiveFilesCantZipMixedValidFile")
  } else {
    t.Logf("%s was created!", archivePath2)
  }
  // ---------------------------------------
  // Pass empty array
  // ---------------------------------------
  var mt []string
  failedArchivePath = archiveFiles(mt)
  if failedArchivePath != "" {
    t.Error("FAILED: archiveFilesOfEmptyArrayDoesNotReturnNil")
  }

  // ---------------------------------------
  // Cleanup
  // ---------------------------------------
  HandleError(t, removeTempFileTempDirError, os.Remove(file.Name()))
  HandleError(t, removeTempFileTempDirError, os.Remove(archivePath1))
  HandleError(t, removeTempFileTempDirError, os.Remove(archivePath2))
  t.Logf("%s was removed!", file.Name())
  t.Logf("%s was removed!", archivePath1)
  t.Logf("%s was removed!", archivePath2)
}

func TestUnarchiveFiles(t *testing.T) {
  // ---------------------------------------
  // Pass invalid zip file
  // ---------------------------------------
  _, err := unarchiveFiles("1")
  if err == nil {
    t.Error("FAILED: unarchiveFilesUnzipsInvalidZip")
  }
  // ---------------------------------------
  // Pass non-existing zip file
  // ---------------------------------------
  _, err = unarchiveFiles(path.Join(tmpDir(), "file.zip"))
  if err == nil {
    t.Error("FAILED: unarchiveFilesUnzipsNonExistingZip")
  }

  // ---------------------------------------
  // Pass existing zip file
  // ---------------------------------------
  zipFile := CreateZipFile(t)
  unzippedFile, err := unarchiveFiles(zipFile)
  if err != nil {
    t.Error("FAILED: unarchiveFilesCanNotUnzipZip")
  }

  zipFileInfo := pathInfo(zipFile)
  if zipFileInfo.Exists {
    t.Error("FAILED: unarchiveFilesDoesNotRemoveZipAfterUnzip")
  }

  // ---------------------------------------
  // Cleanup
  // ---------------------------------------
  // The case where the unzip is successful, unarchiveFiles deletes the zip after unzipping occurs,
  // so we don't have to clean it up
  HandleError(t, removeTempFileTempDirError, os.RemoveAll(unzippedFile))
  t.Logf("%s and its contents was removed!", unzippedFile)
}
