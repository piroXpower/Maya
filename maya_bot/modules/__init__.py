from maya_bot import logger, CONFIG

NO_LOAD_MODULES = CONFIG["advanced"]["not_load_this_modules"]


def list_all_modules():
    from os.path import dirname, basename, isfile
    import glob

    modules = []
    mod_paths = glob.glob(dirname(__file__) + "/*.py")
    all_modules = [
        basename(f)[:-3]
        for f in mod_paths
        if isfile(f) and f.endswith(".py") and not f.endswith("__init__.py")
    ]
    for module in all_modules:
        if module not in NO_LOAD_MODULES:
            modules.append(module)
    return modules


ALL_MODULES = sorted(list_all_modules())
logger.info("Modules to load: %s", str(ALL_MODULES))
__all__ = ALL_MODULES + ["ALL_MODULES"]
