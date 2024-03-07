package mkfgofiles

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	testfilename string
	moddir       string
	testdatadir  string
)

func TestMain(m *testing.M) {
	setup()

	di, errstat := os.Stat(testdatadir)
	if errstat != nil {

		fmt.Println("Test main error: ", errstat)
		panic(errstat)
	}

	if !di.IsDir() {

		fmt.Println("Test main error: ", errors.New("Testdata directory is not a directoy."))
	}

	os.Exit(m.Run())
}

func setup() {

	_, testfilename, _, _ = runtime.Caller(0)
	moddir = filepath.Dir(testfilename)
	testdatadir = filepath.Join(moddir, "testdata")
}

func TestDirectoryExists_Direxists(t *testing.T) {

	_, testfile, _, _ := runtime.Caller(0)
	adir := filepath.Dir(testfile)

	fi, errstat := os.Stat(adir)
	if errstat != nil {

		if os.IsNotExist(errstat) {

			t.Fatalf("%s test failed.  Internal test error.  How the hell did you get here?", t.Name())
		} else {

			t.Fatalf("%s test failed.  Internal test error.  Unexpected test error \"%s\"", t.Name(), errstat)
		}
	}

	if fi.IsDir() {

		de, errde := DirectoryExists(adir)
		if errde != nil {

			t.Fatalf("%s test failed.  Unexpected error \"%s\" returned.", t.Name(), errde)
		}

		if !de {

			t.Fatalf("%s test failed.  Expected true.", t.Name())
		}

	} else {

		t.Fatalf("%s test failed.  Internal test error.  Again, how the hell did you get here?", t.Name())
	}
}

func TestDirectoryExists_DirDoesntExist(t *testing.T) {

	_, testfile, _, _ := runtime.Caller(0)
	adir := filepath.Join(filepath.Dir(testfile), "ishouldnotexist")

	_, errstat := os.Stat(adir)
	if errstat != nil {

		if os.IsNotExist(errstat) {

			de, errde := DirectoryExists(adir)
			if errde != nil {

				t.Fatalf("%s test failed.  Unexpected error \"%s\" returned.", t.Name(), errde)
			}

			if de {

				t.Fatalf("%s test failed.  Expected false.", t.Name())
			}

		} else {

			t.Fatalf("%s test failed.  Internal test error.  Unexpected test error \"%s\"", t.Name(), errstat)
		}
	} else {

		t.Fatalf("%s test failed.  Internal test error.  No error returned when one was expected.", t.Name())
	}
}

func TestRegularfileExists_FileExists(t *testing.T) {

	_, testfile, _, _ := runtime.Caller(0)

	fi, errstat := os.Stat(testfile)

	if errstat != nil {

		if os.IsNotExist(errstat) {

			t.Fatalf("%s - internal test error.  Wut?", t.Name())
		} else {

			t.Fatalf("%s - internal test error.  Unexpected error \"%s\".", t.Name(), errstat)
		}
	}

	if fi.Mode().IsDir() {

		t.Fatalf("%s - internal test error.  It's a directory.", t.Name())
	}

	if fi.Mode().IsRegular() {

		ex, errex := RegularfileExists(testfile)
		if errex != nil {

			t.Fatalf("%s - unexpected error \"%s\"", t.Name(), errex)
		}

		if !ex {

			t.Fatalf("%s - Expected true.", t.Name())
		}

	} else {

		t.Fatalf("%s - internal test error.  It's not a regular file.", t.Name())
	}
}

func TestRegularfileExists_FileDoesntExist(t *testing.T) {

	_, testfile, _, _ := runtime.Caller(0)

	thefile := filepath.Join(filepath.Dir(testfile), "ishouldnotexist.file")

	_, errstat := os.Stat(thefile)

	if errstat != nil {

		if os.IsNotExist(errstat) {

			ex, errex := RegularfileExists(thefile)
			if errex != nil {

				t.Fatalf("%s failed.  Unexpected error encountered \"%s\"", t.Name(), errex)
			}

			if ex {

				t.Fatalf("%s failed.  Expected false.", t.Name())
			}

		} else {

			t.Fatalf("%s - internal test error.  Unexpected error \"%s\".", t.Name(), errstat)
		}
	}
}

func TestRemoveDirectoryWithContents(t *testing.T) {

	testbasedir := filepath.Join(testdatadir, "removecontentstest")

	_, errtd := os.Stat(testbasedir)
	if errtd != nil {

		if !os.IsNotExist(errtd) {

			t.Fatalf("%s failed.  Unexpected error \"%s\" when checking if directory exists.", t.Name(), errtd)
		}
	} else {

		t.Fatalf("%s failed.  Test directory \"%s\" exists and should not at this point.  Remove manually before running test.", t.Name(), testbasedir)
	}

	errmkdir := os.MkdirAll(filepath.Join(testbasedir, "directory", "structure"), 0777)
	if errmkdir != nil {

		t.Fatalf("%s failed.  Could not create test directory.", t.Name())
	}

	s := "some test data"
	b := []byte(s)

	errwrt := os.WriteFile(filepath.Join(testbasedir, "testfile.txt"), b, 0666)
	if errwrt != nil {

		t.Fatalf("%s failed.  Writing internal test file base dir.", t.Name())
	}

	errwrt = os.WriteFile(filepath.Join(testbasedir, "directory", "testfile.txt"), b, 0666)
	if errwrt != nil {

		t.Fatalf("%s failed.  Writing internal test file directory dir.", t.Name())
	}

	errwrt = os.WriteFile(filepath.Join(testbasedir, "directory", "structure", "testfile.txt"), b, 0666)
	if errwrt != nil {

		t.Fatalf("%s failed.  Writing internal test file directory/structure dir.", t.Name())
	}

	errrdc := RemoveDirectoryWithContents(testbasedir)
	if errrdc != nil {

		t.Fatalf("%s failed.  Exuxpected error: \"%s\".", t.Name(), errrdc)
	}

	_, errtd = os.Stat(testbasedir)
	if errtd != nil {

		if !os.IsNotExist(errtd) {

			t.Fatalf("%s failed.  Unexpected error \"%s\" when checking if directory exists after call to remove.", t.Name(), errtd)
		}
	} else {

		t.Fatalf("%s failed.  Test directory exists and should not - failed to remove.", t.Name())
	}

}

func TestRemoveDirectoryWithContentsDirdoentExist(t *testing.T) {

	testdir := filepath.Join(testdatadir, "ishouldntexist")

	_, errstat := os.Stat(testdir)
	if errstat != nil {

		if !os.IsNotExist(errstat) {

			t.Fatalf("%s failed.  Unexpected error in stat \"%s\".", t.Name(), errstat)
		}
	} else {

		t.Fatalf("%s failed.  Directory  \"%s\" should not exist for this test.", t.Name(), testdir)
	}

	expected := fmt.Sprintf("open %s: The system cannot find the file specified.", testdir)

	errne := RemoveDirectoryWithContents(testdir)
	if errne != nil {

		if errne.Error() != expected {

			t.Fatalf("%s failed.  Expected error \"%s\".  Got \"%s\".", t.Name(), expected, errne)
		}

	} else {

		t.Fatalf("%s failed.  Expected error \"%s\".  No error thrown.", t.Name(), expected)
	}

}

func TestRemoveDirectoryWithContentsDirIsAFile(t *testing.T) {

	testdir := filepath.Join(testdatadir, "testfile")

	fi, errstat := os.Stat(testdir)
	if errstat != nil {

		t.Fatalf("%s failed.  Unexpected error \"%s\".", t.Name(), errstat)
	}

	if !fi.Mode().IsRegular() {

		t.Fatalf("%s failed.  Test file \"%s\" should be a regular file.", t.Name(), testdatadir)
	}

	expected := fmt.Sprintf("readdir %s: The system cannot find the path specified.", testdir)

	errremdr := RemoveDirectoryWithContents(testdir)
	if errremdr != nil {

		if errremdr.Error() != expected {

			t.Fatalf("%s failed.  Expected error \"%s\".  Got \"%s\".", t.Name(), expected, errremdr)
		}
	} else {

		t.Fatalf("%s failed.  No error produced.", t.Name())
	}
}
