package copycmd

import (
        "io"
        "os"
)

func CopyFile(source, dest string) error {
        srcFile, err := os.Open(source)
        if err != nil {
                return err
        }
        defer srcFile.Close()

        destFile, err := os.Create(dest)
        if err != nil {
                return err
        }
        defer destFile.Close()

        _, err = io.Copy(destFile, srcFile)
        if err != nil {
                return err
        }

        return nil
}


