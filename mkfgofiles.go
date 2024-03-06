package mkffiles

import (
	"os"
	"path/filepath"
)

func DirectoryExists(adirpath string) (bool, error) {

	di, errexist := os.Stat(adirpath)

	if errexist != nil {

		if os.IsNotExist(errexist) {

			return false, nil
		}

		return false, errexist
	}

	return di.IsDir(), nil
}

func RegularfileExists(afilepath string) (bool, error) {

	fi, errstat := os.Stat(afilepath)

	if errstat != nil {

		if os.IsNotExist(errstat) {

			return false, nil
		}

		return false, errstat
	}

	if fi.IsDir() {

		return false, nil
	}

	return fi.Mode().IsRegular(), nil
}

func RemoveDirectoryWithContents(adirpath string) error {

	d, erropen := os.Open(adirpath)

	if erropen != nil {

		return erropen
	}

	dirnames, errrddr := d.Readdirnames(-1)
	if errrddr != nil {

		return errrddr
	}

	if errdclose := d.Close(); errdclose != nil {

		return errdclose
	}

	for _, dirname := range dirnames {

		errrm := os.RemoveAll(filepath.Join(adirpath, dirname))
		if errrm != nil {

			return errrm
		}
	}

	return os.RemoveAll(adirpath)
}
