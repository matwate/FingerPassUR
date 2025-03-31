#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <glib.h>
#include <glib-unix.h>
#include <libfprint-2/fprint.h>
#include <getopt.h>  

typedef struct _UserData {
  GMainLoop     *loop;
  GCancellable  *cancellable;
  unsigned int  sigint_handler;
  int           ret_value;
  FpPrint       *loaded;
} UserData;

typedef struct _ThreadData {
  GArray  *array;
  guint   start;
  guint   end;
  gint    thread_id;
} ThreadData;

FpDevice *
discover_device(GPtrArray *devices) {
  FpDevice *dev;
  int i;

  if(!devices->len)
    return NULL;

  if(devices->len == 1)
    i = 0;
  else
    return NULL;

  dev = g_ptr_array_index(devices, i);
  g_print ("[%d] %s (%s) - driver %s\n", i,
           fp_device_get_device_id(dev),
           fp_device_get_name(dev),
           fp_device_get_driver(dev));

  return dev;
}

static void
on_enroll_progress(FpDevice *device,
                   gint     completed_stages,
                   FpPrint  *print,
                   gpointer user_data,
                   GError   *error) {
  if(error) {
    g_printerr("Enroll stage %d of %d failed. with error %s\n", 
               completed_stages,
               fp_device_get_nr_enroll_stages(device),
               error->message);
    return;
  }

  if(fp_print_get_image(print)) {
    g_print("\nRead print image! ");
  }

  g_print("Enroll stage %d of %d passed!\n", 
          completed_stages,
          fp_device_get_nr_enroll_stages(device));
}

static void
on_device_closed(FpDevice *dev, GAsyncResult *res, void *user_data) {
  UserData *u_data = user_data;

  g_autoptr(GError) error = NULL;
  fp_device_close_finish(dev, res, &error);

  if(error) {
    g_printerr("Error during device closing %s\n", error->message);
  }

  g_main_loop_quit(u_data->loop);
}

guchar**
get_print_str_repr(FpPrint *print, guchar **out_data, gsize *len) {
  g_autoptr(GError) error = NULL;

  fp_print_serialize(print, out_data, len, &error);
  if(error) {
    g_printerr("Error serializing data %s\n", error->message);
    return NULL;
  }

  return out_data;
}

static void
on_enroll_completed(FpDevice *dev, GAsyncResult *res, void *user_data) {
  UserData *u_data = user_data;

  g_autoptr(FpPrint) print = NULL;
  g_autoptr(GError) error = NULL;

  print = fp_device_enroll_finish(dev, res, &error);

  if(!error) { 
    u_data->ret_value = EXIT_SUCCESS;

    gsize len = 0;
    guchar *data = NULL;
    get_print_str_repr(print, &data, &len);
    
    FILE *f = fopen("print.bin", "wb");
    if(f) {
      fwrite(data, 1, len, f);
      fclose(f);
      g_print("Enrollment complete. Saved to print.bin\n");
    } else {
      g_printerr("Failed to open print.bin for writing!\n");
    }
    g_free(data);

  } else {
    g_printerr("Enrollment failed with error %s\n", error->message);
  }

  fp_device_close(dev, NULL, (GAsyncReadyCallback) on_device_closed, u_data);
}

static void
on_device_opened(FpDevice *dev, GAsyncResult *res, void *user_data) {
  UserData *u_data = user_data;
  FpPrint *print_template;

  g_autoptr(GError) error = NULL;

  if(!fp_device_open_finish(dev, res, &error)) {
    g_printerr("Failed to open device %s\n", error->message);
    g_main_loop_quit(u_data->loop);
    return;
  }

  g_print("Opened device for ENROLL\n");
  g_print("# Scan times: %d\n", fp_device_get_nr_enroll_stages(dev));

  print_template = fp_print_new(dev);

  fp_device_enroll(dev,
                   print_template,
                   u_data->cancellable,
                   on_enroll_progress,
                   NULL, // progress_data
                   NULL, // GDestroyNotify
                   (GAsyncReadyCallback) on_enroll_completed,
                   user_data);
}

FpPrint *
deserialize_print(const char *file_path) {
  FILE *fp = fopen(file_path, "rb");

  if(!fp) {
    g_printerr("Error reading print from %s\n", file_path);
    return NULL;
  }

  fseek(fp, 0, SEEK_END);
  gsize len = ftell(fp);
  rewind(fp);

  guchar *data = g_malloc(len);
  if(fread(data, 1, len, fp) != len) {
    g_printerr("Error reading print data\n");
    g_free(data);
    fclose(fp);
    return NULL;
  }

  g_autoptr(GError) error = NULL;
  FpPrint *print = fp_print_deserialize(data, len, &error);
  g_free(data);
  fclose(fp);


  if(!print) {
    g_printerr("Deserialization failed: %s\n", (error ? error->message : ""));
    return NULL;
  }

  return print;
}

static void
device_verify_callback_finish(FpDevice *dev, GAsyncResult *res, void *user_data) {
  gboolean match;
  g_autoptr(GError) error = NULL;
  FpPrint *scanned_print = NULL;

  fp_device_verify_finish(dev, res, &match, &scanned_print, &error);

  if(error) {
    g_printerr("Error verifying: %s\n", error->message);
  } else if(match) {
    g_print("Matching fingerprints!\n");
  } else {
    g_print("No match.\n");
  }

  if(scanned_print)
    g_object_unref(scanned_print);


  fp_device_close(dev, NULL, (GAsyncReadyCallback) on_device_closed, user_data);

}

static void
on_opened_for_verify(FpDevice *dev,
                     GAsyncResult *res,
                     void *user_data) {

  UserData *u_data = user_data;
  g_autoptr(GError) error = NULL;

  if(!fp_device_open_finish(dev, res, &error)) {
    g_printerr("Failed to open device for verification: %s\n", error->message);
    g_main_loop_quit(u_data->loop);
    return;
  }

  g_print("Opened device for VERIFY\n");

  if(!u_data->loaded) {
    g_printerr("No valid fingerprint template loaded.\n");
    g_main_loop_quit(u_data->loop);
    return;
  }

  fp_device_verify(dev,
                   u_data->loaded,
                   NULL,  // cancellable
                   NULL,  // match_cb
                   NULL,  // match_data
                   NULL,  // match_destroy
                   (GAsyncReadyCallback) device_verify_callback_finish,
                   user_data);
}

static void
on_capture_image_complete(FpDevice *device,
                          GAsyncResult *res,
                          void *user_data) {

  UserData *u_data = user_data;
  g_autoptr(GError) error = NULL;

  FpImage *img = fp_device_capture_finish(device,
                                          res,
                                          &error);

  if(error) {
    g_printerr("Failed finishing the capture %s\n", error->message);
  }
  FILE *f = fopen("img.pgm", "wb");

  gsize len = 0;
  guint width = fp_image_get_width(img);
  guint height = fp_image_get_height(img);

  const guchar *imgdata = fp_image_get_data(img, &len);
  
  fprintf(f, "P5\n%d %d\n255\n", width, height);
  fwrite(imgdata, 1, len, f);
  fclose(f);

  fp_device_close(device, NULL, (GAsyncReadyCallback) on_device_closed, u_data);
  
}

static void
on_opened_for_image(FpDevice *device,
                    GAsyncResult *res,
                    void *user_data) {
  UserData *u_data = user_data;
  g_autoptr(GError) error = NULL;

  fp_device_capture(device,
                    TRUE,
                    u_data->cancellable,
                    (GAsyncReadyCallback) on_capture_image_complete,
                    user_data);

}

static gboolean
sigint_cb(void *user_data) {
  UserData *enroll_data = user_data;
  g_cancellable_cancel(enroll_data->cancellable);
  return G_SOURCE_CONTINUE;
}

int
main(int argc, char *argv[]) {
  gboolean do_enroll = FALSE;
  gboolean do_image = FALSE;
  gchar *verify_file = NULL;

  static struct option long_options[] = {
    {"enroll",      no_argument,       0, 'e'},
    {"verify-with", required_argument, 0, 'v'},
    {"create-image", no_argument,      0, 'i'},
    {0, 0, 0, 0}
  };

  while (1) {
    int opt_index = 0;
    int c = getopt_long(argc, argv, "evi:", long_options, &opt_index);
    if(c == -1)
      break;

    switch(c) {
      case 'e':
        do_enroll = TRUE;
        break;
      case 'v':
        verify_file = g_strdup(optarg);
        break;
      case 'i':
        do_image = TRUE;
        break;
      default:
        g_printerr("Usage: %s [--enroll] [--verify-with file] [--create-image]\n", argv[0]);
        return EXIT_FAILURE;
    }
  }

  if(!do_enroll && !verify_file && !do_image) {
    g_printerr("No operation specified. Use --enroll or --verify-with <file>\n");
    return EXIT_FAILURE;
  }

  g_autoptr(FpContext) ctx = fp_context_new();
  fp_context_enumerate(ctx); 

  FpDevice *dev = discover_device(fp_context_get_devices(ctx));
  if(!dev) {
    g_printerr("Could not find any device\n");
    return EXIT_FAILURE;
  }

  UserData *user_data = g_new0(UserData, 1);
  user_data->ret_value = EXIT_FAILURE;
  user_data->loop = g_main_loop_new(NULL, FALSE);
  user_data->cancellable = g_cancellable_new();
  user_data->sigint_handler = g_unix_signal_add_full(
                                G_PRIORITY_HIGH,
                                SIGINT,
                                sigint_cb,
                                user_data,
                                NULL);

  if(do_enroll) {
    fp_device_open(dev,
                   user_data->cancellable,
                   (GAsyncReadyCallback) on_device_opened,
                   user_data);
  } if (do_image) {
    fp_device_open(dev,
                   user_data->cancellable,
                   (GAsyncReadyCallback) on_opened_for_image,
                   user_data);
  } else {
    // verify_file must be set if we get here
    g_print("Loading template from: %s\n", verify_file);
    user_data->loaded = deserialize_print(verify_file);
    // Now open for VERIFY
    fp_device_open(dev,
                   user_data->cancellable,
                   (GAsyncReadyCallback) on_opened_for_verify,
                   user_data);
  }

  g_main_loop_run(user_data->loop);

  if(verify_file)
    g_free(verify_file);

  g_main_loop_unref(user_data->loop);
  g_free(user_data);

  return EXIT_SUCCESS;
}

