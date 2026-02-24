init() {
    level._dvars = SpawnStruct();
    level._dvars.enabled = "plutoplugin_enabled";
    level._dvars.inDvar  = "plutoplugin_in";
    level._dvars.outDvar = "plutoplugin_out";
    level._dvars.reset   = "";

    level._commnads = [];
    level._command_prefix  = "!";

    scripts\mp\_plutoplugin_utils::SetDvarIfUnitialized(level._dvars.enabled, 1);
    scripts\mp\_plutoplugin_utils::SetDvarIfUnitialized(level._dvars.inDvar, "");
    scripts\mp\_plutoplugin_utils::SetDvarIfUnitialized(level._dvars.inDvar, "");

    EnableDvarChangedNotify(level._dvars.enabled);
    EnableDvarChangedNotify(level._dvars.inDvar);
}

