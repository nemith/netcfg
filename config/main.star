load("schema/site.proto", "Site")

register_object(Site, glob(["site/**/*.star"]))